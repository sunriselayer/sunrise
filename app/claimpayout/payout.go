// Package claimpayout checks a Keplr ADR-036 signature against claims.json
// and records an external USDC or Cosmos payout.
//
// The signature proves that the sunrise account owner asked for a specific
// amount, asset, and destination. It does not make that account an Edge owner.
// RISE rows are rejected here because they are paid later by an Edge claim.
//
// This package does not hold hot-wallet private keys and does not broadcast.
// Reserve a nonce, send from the hot wallet outside this process, then record
// the transaction hash. A failed send releases the reservation. Caps are
// supplied by the operator for each asset; a missing cap rejects the payout.
package claimpayout

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cometbft/cometbft/crypto/secp256k1"
)

const (
	payoutUSDC   = "usdc"
	payoutCosmos = "cosmos"
	payoutEdge   = "edge"
)

var evmAddress = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

// Authorization is the message a claimant signs with Keplr.
// Canonical JSON field order is amount, asset, claimant, destination,
// ledger_sha256, nonce.
type Authorization struct {
	LedgerSHA256 string `json:"ledger_sha256"`
	Claimant     string `json:"claimant"`
	Asset        string `json:"asset"`
	Amount       string `json:"amount"`
	Destination  string `json:"destination"`
	Nonce        uint64 `json:"nonce"`
}

type canonicalAuthorization struct {
	Amount       string `json:"amount"`
	Asset        string `json:"asset"`
	Claimant     string `json:"claimant"`
	Destination  string `json:"destination"`
	LedgerSHA256 string `json:"ledger_sha256"`
	Nonce        uint64 `json:"nonce"`
}

type signDoc struct {
	AccountNumber string    `json:"account_number"`
	ChainID       string    `json:"chain_id"`
	Fee           signFee   `json:"fee"`
	Memo          string    `json:"memo"`
	Msgs          []signMsg `json:"msgs"`
	Sequence      string    `json:"sequence"`
}

type signFee struct {
	Amount []struct{} `json:"amount"`
	Gas    string     `json:"gas"`
}

type signMsg struct {
	Type  string        `json:"type"`
	Value signDataValue `json:"value"`
}

type signDataValue struct {
	Data   string `json:"data"`
	Signer string `json:"signer"`
}

// Policy is the operator limit for one asset, in that asset's base units.
type Policy struct {
	MaxPerTx math.Int
	MaxDaily math.Int
}

// Claim is the payable subset of a claims.json row.
type Claim struct {
	Owner       string `json:"owner"`
	Asset       string `json:"asset"`
	Amount      string `json:"amount"`
	ClaimableAt int64  `json:"claimable_at"`
	Payout      string `json:"payout"`
}

type claimsDocument struct {
	Claims []Claim `json:"claims"`
}

// Reservation is a nonce held while the hot wallet send is in progress.
type Reservation struct {
	Nonce        uint64 `json:"nonce"`
	LedgerSHA256 string `json:"ledger_sha256"`
	Claimant     string `json:"claimant"`
	Asset        string `json:"asset"`
	Amount       string `json:"amount"`
	Destination  string `json:"destination"`
	Payout       string `json:"payout"`
	ReservedAt   string `json:"reserved_at"`
}

// Consumption is a payout that already has a transaction hash.
type Consumption struct {
	Reservation
	TxHash      string `json:"tx_hash"`
	CompletedAt string `json:"completed_at"`
}

type storeFile struct {
	Reservations []Reservation `json:"reservations"`
	Consumptions []Consumption `json:"consumptions"`
}

// Sender broadcasts one already-checked payout. Implementations live outside
// this repository so private keys are not imported here.
type Sender interface {
	Send(payout, asset, destination, amount string) (txHash string, err error)
}

func CanonicalBytes(auth Authorization) ([]byte, error) {
	encoded, err := json.Marshal(canonicalAuthorization{
		Amount:       auth.Amount,
		Asset:        auth.Asset,
		Claimant:     auth.Claimant,
		Destination:  auth.Destination,
		LedgerSHA256: auth.LedgerSHA256,
		Nonce:        auth.Nonce,
	})
	if err != nil {
		return nil, fmt.Errorf("CanonicalBytes: nonce %d claimant %s: %w", auth.Nonce, auth.Claimant, err)
	}
	return encoded, nil
}

func SignDocBytes(canonical []byte, claimant string) ([]byte, error) {
	encoded, err := json.Marshal(signDoc{
		AccountNumber: "0",
		ChainID:       "",
		Fee: signFee{
			Amount: []struct{}{},
			Gas:    "0",
		},
		Memo: "",
		Msgs: []signMsg{{
			Type: "sign/MsgSignData",
			Value: signDataValue{
				Data:   base64.StdEncoding.EncodeToString(canonical),
				Signer: claimant,
			},
		}},
		Sequence: "0",
	})
	if err != nil {
		return nil, fmt.Errorf("SignDocBytes: claimant %s: %w", claimant, err)
	}
	return encoded, nil
}

func AddressFromPubKey(pub []byte) (string, error) {
	if len(pub) != secp256k1.PubKeySize {
		return "", fmt.Errorf("AddressFromPubKey: pubkey length %d", len(pub))
	}
	return sdk.AccAddress(secp256k1.PubKey(pub).Address()).String(), nil
}

func VerifySignature(pub, signDocBytes, sig []byte) error {
	if len(pub) != secp256k1.PubKeySize {
		return fmt.Errorf("VerifySignature: pubkey length %d", len(pub))
	}
	if len(sig) != 64 {
		return fmt.Errorf("VerifySignature: signature length %d", len(sig))
	}
	if !secp256k1.PubKey(pub).VerifySignature(signDocBytes, sig) {
		return fmt.Errorf("VerifySignature: signature does not match the pubkey")
	}
	return nil
}

// Check confirms the signature and reports the remaining balance.
// An empty storePath skips recorded payouts.
func Check(ledgerPath, claimsPath, storePath string, auth Authorization, pub, sig []byte, now time.Time) (math.Int, string, error) {
	ledgerHash, err := FileSHA256(ledgerPath)
	if err != nil {
		return math.Int{}, "", err
	}
	if !strings.EqualFold(auth.LedgerSHA256, ledgerHash) {
		return math.Int{}, "", fmt.Errorf("Check: ledger hash %s does not match %s", auth.LedgerSHA256, ledgerHash)
	}
	derived, err := AddressFromPubKey(pub)
	if err != nil {
		return math.Int{}, "", fmt.Errorf("Check: %w", err)
	}
	if derived != auth.Claimant {
		return math.Int{}, "", fmt.Errorf("Check: pubkey address %s does not match claimant %s", derived, auth.Claimant)
	}
	canonical, err := CanonicalBytes(auth)
	if err != nil {
		return math.Int{}, "", err
	}
	signDocBytes, err := SignDocBytes(canonical, auth.Claimant)
	if err != nil {
		return math.Int{}, "", err
	}
	if err := VerifySignature(pub, signDocBytes, sig); err != nil {
		return math.Int{}, "", fmt.Errorf("Check: claimant %s: %w", auth.Claimant, err)
	}
	claims, err := LoadClaims(claimsPath)
	if err != nil {
		return math.Int{}, "", err
	}
	store := storeFile{Reservations: []Reservation{}, Consumptions: []Consumption{}}
	if storePath != "" {
		loaded, err := loadStore(storePath)
		if err != nil {
			return math.Int{}, "", err
		}
		store = loaded
	}
	return Available(claims, store, auth.Claimant, auth.Asset, now)
}

func FileSHA256(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("FileSHA256: read %s: %w", path, err)
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

func LoadClaims(path string) ([]Claim, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("LoadClaims: read %s: %w", path, err)
	}
	var doc claimsDocument
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("LoadClaims: unmarshal %s: %w", path, err)
	}
	if doc.Claims == nil {
		doc.Claims = []Claim{}
	}
	return doc.Claims, nil
}

// Available is the claimable amount minus reservations and completed payouts.
// Records for a different ledger file still reduce the balance, so publishing
// a replacement claims file does not make the same balance payable again.
func Available(claims []Claim, store storeFile, claimant, asset string, now time.Time) (math.Int, string, error) {
	total := math.ZeroInt()
	payout := ""
	for _, claim := range claims {
		if claim.Owner != claimant || claim.Asset != asset || claim.ClaimableAt > now.Unix() {
			continue
		}
		if payout == "" {
			payout = claim.Payout
		} else if payout != claim.Payout {
			return math.Int{}, "", fmt.Errorf("Available: claimant %s asset %s has payouts %s and %s", claimant, asset, payout, claim.Payout)
		}
		amount, err := parseAmount(claim.Amount)
		if err != nil {
			return math.Int{}, "", fmt.Errorf("Available: claimant %s asset %s: %w", claimant, asset, err)
		}
		total = total.Add(amount)
	}
	spent, err := spentAmount(store, claimant, asset)
	if err != nil {
		return math.Int{}, "", err
	}
	if spent.GT(total) {
		return math.Int{}, "", fmt.Errorf("Available: claimant %s asset %s recorded %s exceeds claim total %s", claimant, asset, spent, total)
	}
	return total.Sub(spent), payout, nil
}

func spentAmount(store storeFile, claimant, asset string) (math.Int, error) {
	spent := math.ZeroInt()
	for _, item := range store.Reservations {
		if item.Claimant != claimant || item.Asset != asset {
			continue
		}
		amount, err := parseAmount(item.Amount)
		if err != nil {
			return math.Int{}, fmt.Errorf("spentAmount: reservation nonce %d: %w", item.Nonce, err)
		}
		spent = spent.Add(amount)
	}
	for _, item := range store.Consumptions {
		if item.Claimant != claimant || item.Asset != asset {
			continue
		}
		amount, err := parseAmount(item.Amount)
		if err != nil {
			return math.Int{}, fmt.Errorf("spentAmount: consumption nonce %d: %w", item.Nonce, err)
		}
		spent = spent.Add(amount)
	}
	return spent, nil
}

// Reserve checks the signature, the remaining balance, and the caps, then
// holds the nonce until Complete or Release.
func Reserve(storePath, ledgerPath, claimsPath string, auth Authorization, pub, sig []byte, policy Policy, now time.Time) error {
	if err := validatePolicy(policy, auth); err != nil {
		return err
	}
	amount, err := parseAmount(auth.Amount)
	if err != nil {
		return fmt.Errorf("Reserve: nonce %d: %w", auth.Nonce, err)
	}
	ledgerHash, err := FileSHA256(ledgerPath)
	if err != nil {
		return err
	}
	if !strings.EqualFold(auth.LedgerSHA256, ledgerHash) {
		return fmt.Errorf("Reserve: nonce %d ledger hash %s does not match %s", auth.Nonce, auth.LedgerSHA256, ledgerHash)
	}
	derived, err := AddressFromPubKey(pub)
	if err != nil {
		return fmt.Errorf("Reserve: nonce %d: %w", auth.Nonce, err)
	}
	if derived != auth.Claimant {
		return fmt.Errorf("Reserve: nonce %d pubkey address %s does not match claimant %s", auth.Nonce, derived, auth.Claimant)
	}
	canonical, err := CanonicalBytes(auth)
	if err != nil {
		return err
	}
	signDocBytes, err := SignDocBytes(canonical, auth.Claimant)
	if err != nil {
		return err
	}
	if err := VerifySignature(pub, signDocBytes, sig); err != nil {
		return fmt.Errorf("Reserve: nonce %d claimant %s: %w", auth.Nonce, auth.Claimant, err)
	}
	claims, err := LoadClaims(claimsPath)
	if err != nil {
		return err
	}
	return mutateStore(storePath, func(store *storeFile) error {
		if err := rejectKnownNonce(*store, auth.Nonce); err != nil {
			return err
		}
		available, payout, err := Available(claims, *store, auth.Claimant, auth.Asset, now)
		if err != nil {
			return err
		}
		if payout == payoutEdge {
			return fmt.Errorf("Reserve: nonce %d asset %s is paid on Edge, not by an external transfer", auth.Nonce, auth.Asset)
		}
		if payout != payoutUSDC && payout != payoutCosmos {
			return fmt.Errorf("Reserve: nonce %d asset %s has no external payout", auth.Nonce, auth.Asset)
		}
		if err := validateDestination(payout, auth.Destination); err != nil {
			return fmt.Errorf("Reserve: nonce %d: %w", auth.Nonce, err)
		}
		if amount.GT(available) {
			return fmt.Errorf("Reserve: nonce %d amount %s exceeds available %s for %s %s", auth.Nonce, amount, available, auth.Claimant, auth.Asset)
		}
		if amount.GT(policy.MaxPerTx) {
			return fmt.Errorf("Reserve: nonce %d amount %s exceeds per-tx cap %s", auth.Nonce, amount, policy.MaxPerTx)
		}
		today, err := spentToday(*store, auth.Asset, now)
		if err != nil {
			return err
		}
		if today.Add(amount).GT(policy.MaxDaily) {
			return fmt.Errorf("Reserve: nonce %d amount %s plus today's %s exceeds daily cap %s for %s", auth.Nonce, amount, today, policy.MaxDaily, auth.Asset)
		}
		store.Reservations = append(store.Reservations, Reservation{
			Nonce:        auth.Nonce,
			LedgerSHA256: strings.ToLower(auth.LedgerSHA256),
			Claimant:     auth.Claimant,
			Asset:        auth.Asset,
			Amount:       amount.String(),
			Destination:  auth.Destination,
			Payout:       payout,
			ReservedAt:   now.UTC().Format(time.RFC3339),
		})
		return nil
	})
}

// Complete stores the hot-wallet transaction hash and consumes the nonce.
func Complete(storePath string, nonce uint64, txHash string, now time.Time) error {
	if txHash == "" {
		return fmt.Errorf("Complete: nonce %d is missing a transaction hash", nonce)
	}
	return mutateStore(storePath, func(store *storeFile) error {
		for i, item := range store.Reservations {
			if item.Nonce != nonce {
				continue
			}
			store.Consumptions = append(store.Consumptions, Consumption{
				Reservation: item,
				TxHash:      txHash,
				CompletedAt: now.UTC().Format(time.RFC3339),
			})
			store.Reservations = append(store.Reservations[:i], store.Reservations[i+1:]...)
			return nil
		}
		return fmt.Errorf("Complete: nonce %d is not reserved", nonce)
	})
}

// Release returns a reserved nonce after the external send fails.
func Release(storePath string, nonce uint64) error {
	return mutateStore(storePath, func(store *storeFile) error {
		for i, item := range store.Reservations {
			if item.Nonce != nonce {
				continue
			}
			store.Reservations = append(store.Reservations[:i], store.Reservations[i+1:]...)
			return nil
		}
		return fmt.Errorf("Release: nonce %d is not reserved", nonce)
	})
}

// Execute reserves the nonce, asks the sender to broadcast, and either records
// the hash or releases the reservation. If the broadcast succeeds and recording
// the hash fails, the reservation stays so the hash can be recorded later.
func Execute(storePath, ledgerPath, claimsPath string, auth Authorization, pub, sig []byte, policy Policy, now time.Time, sender Sender) error {
	if sender == nil {
		return fmt.Errorf("Execute: nonce %d sender is nil", auth.Nonce)
	}
	if err := Reserve(storePath, ledgerPath, claimsPath, auth, pub, sig, policy, now); err != nil {
		return err
	}
	var payout, destination, amount string
	if err := readStore(storePath, func(store storeFile) error {
		for _, item := range store.Reservations {
			if item.Nonce == auth.Nonce {
				payout = item.Payout
				destination = item.Destination
				amount = item.Amount
				return nil
			}
		}
		return fmt.Errorf("Execute: nonce %d disappeared after reserve", auth.Nonce)
	}); err != nil {
		return err
	}
	txHash, err := sender.Send(payout, auth.Asset, destination, amount)
	if err != nil {
		releaseErr := Release(storePath, auth.Nonce)
		return fmt.Errorf("Execute: send nonce %d asset %s to %s: %w (release: %v)", auth.Nonce, auth.Asset, destination, err, releaseErr)
	}
	if err := Complete(storePath, auth.Nonce, txHash, now); err != nil {
		return fmt.Errorf("Execute: broadcast %s for nonce %d succeeded, but recording it failed: %w", txHash, auth.Nonce, err)
	}
	return nil
}

func validatePolicy(policy Policy, auth Authorization) error {
	if policy.MaxPerTx.IsNil() || !policy.MaxPerTx.IsPositive() {
		return fmt.Errorf("validatePolicy: nonce %d asset %s is missing a positive per-tx cap", auth.Nonce, auth.Asset)
	}
	if policy.MaxDaily.IsNil() || !policy.MaxDaily.IsPositive() {
		return fmt.Errorf("validatePolicy: nonce %d asset %s is missing a positive daily cap", auth.Nonce, auth.Asset)
	}
	return nil
}

func validateDestination(payout, destination string) error {
	switch payout {
	case payoutUSDC:
		if !evmAddress.MatchString(destination) {
			return fmt.Errorf("validateDestination: USDC destination %q is not an EVM address", destination)
		}
	case payoutCosmos:
		_, converted, err := bech32.DecodeAndConvert(destination)
		if err != nil || (len(converted) != 20 && len(converted) != 32) {
			return fmt.Errorf("validateDestination: Cosmos destination %q is not a bech32 address", destination)
		}
	default:
		return fmt.Errorf("validateDestination: payout %s is not external", payout)
	}
	return nil
}

func rejectKnownNonce(store storeFile, nonce uint64) error {
	for _, item := range store.Reservations {
		if item.Nonce == nonce {
			return fmt.Errorf("rejectKnownNonce: nonce %d is already reserved", nonce)
		}
	}
	for _, item := range store.Consumptions {
		if item.Nonce == nonce {
			return fmt.Errorf("rejectKnownNonce: nonce %d is already consumed", nonce)
		}
	}
	return nil
}

func spentToday(store storeFile, asset string, now time.Time) (math.Int, error) {
	day := now.UTC().Format("2006-01-02")
	total := math.ZeroInt()
	for _, item := range store.Reservations {
		if item.Asset != asset {
			continue
		}
		amount, err := parseAmount(item.Amount)
		if err != nil {
			return math.Int{}, fmt.Errorf("spentToday: reservation nonce %d: %w", item.Nonce, err)
		}
		total = total.Add(amount)
	}
	for _, item := range store.Consumptions {
		if item.Asset != asset || !sameDay(item.CompletedAt, day) {
			continue
		}
		amount, err := parseAmount(item.Amount)
		if err != nil {
			return math.Int{}, fmt.Errorf("spentToday: consumption nonce %d: %w", item.Nonce, err)
		}
		total = total.Add(amount)
	}
	return total, nil
}

func sameDay(stamp, day string) bool {
	parsed, err := time.Parse(time.RFC3339, stamp)
	if err != nil {
		return false
	}
	return parsed.UTC().Format("2006-01-02") == day
}

func parseAmount(value string) (math.Int, error) {
	parsed, ok := math.NewIntFromString(value)
	if !ok || !parsed.IsPositive() {
		return math.Int{}, fmt.Errorf("parseAmount: %q", value)
	}
	return parsed, nil
}

func DecodeHex(label, value string) ([]byte, error) {
	raw, err := hex.DecodeString(strings.TrimPrefix(value, "0x"))
	if err != nil {
		return nil, fmt.Errorf("DecodeHex: %s %q: %w", label, value, err)
	}
	return raw, nil
}

func readStore(path string, fn func(storeFile) error) error {
	return withLock(path, func() error {
		store, err := loadStore(path)
		if err != nil {
			return err
		}
		return fn(store)
	})
}

func mutateStore(path string, fn func(*storeFile) error) error {
	return withLock(path, func() error {
		store, err := loadStore(path)
		if err != nil {
			return err
		}
		if err := fn(&store); err != nil {
			return err
		}
		if store.Reservations == nil {
			store.Reservations = []Reservation{}
		}
		if store.Consumptions == nil {
			store.Consumptions = []Consumption{}
		}
		encoded, err := json.MarshalIndent(store, "", "  ")
		if err != nil {
			return fmt.Errorf("mutateStore: encode %s: %w", path, err)
		}
		encoded = append(encoded, '\n')
		temp := path + ".tmp"
		if err := os.WriteFile(temp, encoded, 0o600); err != nil {
			return fmt.Errorf("mutateStore: write %s: %w", temp, err)
		}
		if err := os.Rename(temp, path); err != nil {
			return fmt.Errorf("mutateStore: rename %s: %w", path, err)
		}
		return nil
	})
}

func loadStore(path string) (storeFile, error) {
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return storeFile{Reservations: []Reservation{}, Consumptions: []Consumption{}}, nil
	}
	if err != nil {
		return storeFile{}, fmt.Errorf("loadStore: read %s: %w", path, err)
	}
	if len(bytes.TrimSpace(raw)) == 0 {
		return storeFile{Reservations: []Reservation{}, Consumptions: []Consumption{}}, nil
	}
	var store storeFile
	if err := json.Unmarshal(raw, &store); err != nil {
		return storeFile{}, fmt.Errorf("loadStore: unmarshal %s: %w", path, err)
	}
	if store.Reservations == nil {
		store.Reservations = []Reservation{}
	}
	if store.Consumptions == nil {
		store.Consumptions = []Consumption{}
	}
	return store, nil
}
