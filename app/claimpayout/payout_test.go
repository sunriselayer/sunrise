package claimpayout

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("sunrise", "sunrisepub")
	os.Exit(m.Run())
}

func TestCanonicalAuthorizationIsStable(t *testing.T) {
	encoded, err := CanonicalBytes(Authorization{
		LedgerSHA256: "abc",
		Claimant:     "sunrise1example",
		Asset:        "usdn",
		Amount:       "5",
		Destination:  "0x0000000000000000000000000000000000000001",
		Nonce:        7,
	})
	require.NoError(t, err)
	require.Equal(t, `{"amount":"5","asset":"usdn","claimant":"sunrise1example","destination":"0x0000000000000000000000000000000000000001","ledger_sha256":"abc","nonce":7}`, string(encoded))
}

func TestReserveCompleteAndReplay(t *testing.T) {
	now := time.Date(2026, 10, 6, 0, 0, 0, 0, time.UTC)
	claimant, pub, sig, auth, ledgerPath, claimsPath, priv := signedFixture(t, now, "4", "usdn", "0x0000000000000000000000000000000000000001", 1)
	storePath := filepath.Join(t.TempDir(), "payouts.json")
	policy := Policy{MaxPerTx: math.NewInt(10), MaxDaily: math.NewInt(10)}

	require.NoError(t, Reserve(storePath, ledgerPath, claimsPath, auth, pub, sig, policy, now))
	require.Error(t, Reserve(storePath, ledgerPath, claimsPath, auth, pub, sig, policy, now))

	claims, err := LoadClaims(claimsPath)
	require.NoError(t, err)
	store := mustStore(t, storePath)
	available, payout, err := Available(claims, store, claimant, "usdn", now)
	require.NoError(t, err)
	require.Equal(t, payoutUSDC, payout)
	require.Equal(t, "6", available.String())

	require.NoError(t, Complete(storePath, 1, "0xhash", now))
	store = mustStore(t, storePath)
	require.Empty(t, store.Reservations)
	require.Equal(t, "0xhash", store.Consumptions[0].TxHash)
	available, _, err = Available(claims, store, claimant, "usdn", now)
	require.NoError(t, err)
	require.Equal(t, "6", available.String())

	auth.Nonce = 2
	auth.Amount = "6"
	pub, sig = signWith(t, priv, auth)
	require.NoError(t, Reserve(storePath, ledgerPath, claimsPath, auth, pub, sig, policy, now))
	auth.Nonce = 3
	auth.Amount = "1"
	pub, sig = signWith(t, priv, auth)
	require.Error(t, Reserve(storePath, ledgerPath, claimsPath, auth, pub, sig, policy, now))
}

func TestReserveRejectsLockedFundsCapsAndEdge(t *testing.T) {
	now := time.Date(2026, 10, 6, 0, 0, 0, 0, time.UTC)
	_, pub, sig, auth, ledgerPath, claimsPath, priv := signedFixture(t, now, "5", "usdn", "0x0000000000000000000000000000000000000001", 4)
	storePath := filepath.Join(t.TempDir(), "payouts.json")

	require.Error(t, Reserve(storePath, ledgerPath, claimsPath, auth, pub, sig, Policy{MaxPerTx: math.NewInt(4), MaxDaily: math.NewInt(10)}, now))
	require.NoError(t, Reserve(storePath, ledgerPath, claimsPath, auth, pub, sig, Policy{MaxPerTx: math.NewInt(5), MaxDaily: math.NewInt(5)}, now))
	require.NoError(t, Release(storePath, 4))

	auth.Nonce = 5
	auth.Amount = "1"
	auth.Asset = "rise"
	auth.Destination = "edge-address"
	pub, sig = signWith(t, priv, auth)
	require.Error(t, Reserve(storePath, ledgerPath, claimsPath, auth, pub, sig, Policy{MaxPerTx: math.NewInt(10), MaxDaily: math.NewInt(10)}, now))

	future := now.Add(-time.Hour)
	destination, err := bech32.ConvertAndEncode("cosmos", []byte("cosmos-destination01"))
	require.NoError(t, err)
	_, pub, sig, auth, ledgerPath, claimsPath, _ = signedFixture(t, future, "9", "ibc/ATOM", destination, 8)
	require.Error(t, Reserve(storePath, ledgerPath, claimsPath, auth, pub, sig, Policy{MaxPerTx: math.NewInt(100), MaxDaily: math.NewInt(100)}, future))
}

func TestExecuteReleasesWhenSendFails(t *testing.T) {
	now := time.Date(2026, 10, 6, 0, 0, 0, 0, time.UTC)
	_, pub, sig, auth, ledgerPath, claimsPath, _ := signedFixture(t, now, "1", "usdn", "0x0000000000000000000000000000000000000001", 9)
	storePath := filepath.Join(t.TempDir(), "payouts.json")
	policy := Policy{MaxPerTx: math.NewInt(10), MaxDaily: math.NewInt(10)}
	err := Execute(storePath, ledgerPath, claimsPath, auth, pub, sig, policy, now, failingSender{})
	require.Error(t, err)
	store := mustStore(t, storePath)
	require.Empty(t, store.Reservations)
	require.Empty(t, store.Consumptions)

	require.NoError(t, Execute(storePath, ledgerPath, claimsPath, auth, pub, sig, policy, now, staticSender{hash: "tx-1"}))
	store = mustStore(t, storePath)
	require.Equal(t, "tx-1", store.Consumptions[0].TxHash)
}

func TestVerifySignatureRejectsADifferentKey(t *testing.T) {
	now := time.Date(2026, 10, 6, 0, 0, 0, 0, time.UTC)
	_, pub, sig, auth, ledgerPath, claimsPath, _ := signedFixture(t, now, "1", "usdn", "0x0000000000000000000000000000000000000001", 11)
	other := secp256k1.GenPrivKey().PubKey().Bytes()
	storePath := filepath.Join(t.TempDir(), "payouts.json")
	policy := Policy{MaxPerTx: math.NewInt(10), MaxDaily: math.NewInt(10)}
	require.Error(t, Reserve(storePath, ledgerPath, claimsPath, auth, other, sig, policy, now))
	require.NotEqual(t, pub, other)
}

type failingSender struct{}

func (failingSender) Send(string, string, string, string) (string, error) {
	return "", os.ErrClosed
}

type staticSender struct {
	hash string
}

func (s staticSender) Send(string, string, string, string) (string, error) {
	return s.hash, nil
}

func signedFixture(t *testing.T, now time.Time, amount, asset, destination string, nonce uint64) (string, []byte, []byte, Authorization, string, string, secp256k1.PrivKey) {
	t.Helper()
	dir := t.TempDir()
	ledgerPath := filepath.Join(dir, "claims.json")
	body := map[string]any{
		"claims": []any{
			map[string]any{"owner": "", "asset": "usdn", "amount": "10", "claimable_at": now.Unix(), "payout": payoutUSDC},
			map[string]any{"owner": "", "asset": "rise", "amount": "10", "claimable_at": now.Unix(), "payout": payoutEdge},
			map[string]any{"owner": "", "asset": "ibc/ATOM", "amount": "9", "claimable_at": now.Unix() + 3600, "payout": payoutCosmos},
		},
	}
	priv := secp256k1.GenPrivKey()
	claimant := sdk.AccAddress(priv.PubKey().Address()).String()
	claims := body["claims"].([]any)
	for _, item := range claims {
		item.(map[string]any)["owner"] = claimant
	}
	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(ledgerPath, encoded, 0o600))
	sum, err := FileSHA256(ledgerPath)
	require.NoError(t, err)
	auth := Authorization{
		LedgerSHA256: sum,
		Claimant:     claimant,
		Asset:        asset,
		Amount:       amount,
		Destination:  destination,
		Nonce:        nonce,
	}
	pub, sig := signWith(t, priv, auth)
	return claimant, pub, sig, auth, ledgerPath, ledgerPath, priv
}

func signWith(t *testing.T, priv secp256k1.PrivKey, auth Authorization) ([]byte, []byte) {
	t.Helper()
	canonical, err := CanonicalBytes(auth)
	require.NoError(t, err)
	signDoc, err := SignDocBytes(canonical, auth.Claimant)
	require.NoError(t, err)
	sig, err := priv.Sign(signDoc)
	require.NoError(t, err)
	return priv.PubKey().Bytes(), sig
}

func mustStore(t *testing.T, path string) storeFile {
	t.Helper()
	var store storeFile
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &store))
	return store
}
