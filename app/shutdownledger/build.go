// Package shutdownledger builds readable holder files from a pre-upgrade state export.
// The export is the source of record. These files can be regenerated after the chain is shut down.
package shutdownledger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	pooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	lockuptypes "github.com/sunriselayer/sunrise/x/lockup/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
	swaptypes "github.com/sunriselayer/sunrise/x/swap/types"
)

// File names written into the output directory.
const (
	BankFile      = "bank.json"
	PositionsFile = "positions.json"
	LockupFile    = "lockup.json"
	InFlightFile  = "in_flight.json"
	StakingFile   = "staking.json"
	ClaimsFile    = "claims.json"
)

// Build reads an application export and writes the holder files.
// Bank balances are copied as stored. Position token amounts are calculated separately
// and are not added to the bank file.
// snapshotTime is the block time at the export height. The export does not record it,
// so claims.json vests lockups at this time.
func Build(inputPath, outputDir string, snapshotTime time.Time) error {
	raw, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("Build: read %s: %w", inputPath, err)
	}
	var doc exportDocument
	if err := json.Unmarshal(raw, &doc); err != nil {
		return fmt.Errorf("Build: unmarshal %s: %w", inputPath, err)
	}
	if doc.AppState == nil {
		return fmt.Errorf("Build: %s has no app_state", inputPath)
	}

	bankFile, positionsFile, lockupFile, inFlightFile, stakingFile, err := buildFiles(doc.AppState)
	if err != nil {
		return err
	}
	claims, err := buildClaims(doc, snapshotTime)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("Build: create %s: %w", outputDir, err)
	}
	files := map[string]any{
		BankFile:      bankFile,
		PositionsFile: positionsFile,
		LockupFile:    lockupFile,
		InFlightFile:  inFlightFile,
		StakingFile:   stakingFile,
		ClaimsFile:    claims,
	}
	for name, value := range files {
		path := filepath.Join(outputDir, name)
		if err := writeJSON(path, value); err != nil {
			return fmt.Errorf("Build: write %s: %w", path, err)
		}
	}
	return nil
}

func writeJSON(path string, value any) error {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	encoded = append(encoded, '\n')
	return os.WriteFile(path, encoded, 0o600)
}

type bankFile struct {
	Balances []banktypes.Balance `json:"balances"`
}

type positionRecord struct {
	PositionID       uint64       `json:"position_id"`
	Owner            string       `json:"owner"`
	PoolID           uint64       `json:"pool_id"`
	LowerTick        int64        `json:"lower_tick"`
	UpperTick        int64        `json:"upper_tick"`
	Liquidity        string       `json:"liquidity"`
	DenomBase        string       `json:"denom_base"`
	DenomQuote       string       `json:"denom_quote"`
	AmountBase       string       `json:"amount_base"`
	AmountQuote      string       `json:"amount_quote"`
	UnclaimedRewards sdk.DecCoins `json:"unclaimed_rewards"`
}

type lockupRecord struct {
	Owner         string    `json:"owner"`
	LockupAddress string    `json:"lockup_address"`
	ID            uint64    `json:"id"`
	Balances      sdk.Coins `json:"balances"`
}

type swapExport struct {
	Incoming []incomingExport `json:"incoming_in_flight_packets"`
	Outgoing []outgoingExport `json:"outgoing_in_flight_packets"`
}

type packetIndexExport struct {
	PortID    string     `json:"port_id"`
	ChannelID string     `json:"channel_id"`
	Sequence  flexUint64 `json:"sequence"`
}

type incomingExport struct {
	Index packetIndexExport `json:"index"`
	Data  []byte            `json:"data"`
}

type outgoingExport struct {
	Index            packetIndexExport `json:"index"`
	AckWaitingIndex  packetIndexExport `json:"ack_waiting_index"`
	RetriesRemaining flexInt64         `json:"retries_remaining"`
}

type inFlightFile struct {
	Incoming []incomingPacketRecord `json:"incoming"`
	Outgoing []outgoingPacketRecord `json:"outgoing"`
}

type incomingPacketRecord struct {
	PortID    string `json:"port_id"`
	ChannelID string `json:"channel_id"`
	Sequence  uint64 `json:"sequence"`
	Denom     string `json:"denom,omitempty"`
	Amount    string `json:"amount,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Receiver  string `json:"receiver,omitempty"`
}

type outgoingPacketRecord struct {
	PortID            string `json:"port_id"`
	ChannelID         string `json:"channel_id"`
	Sequence          uint64 `json:"sequence"`
	AckWaitingPortID  string `json:"ack_waiting_port_id,omitempty"`
	AckWaitingChannel string `json:"ack_waiting_channel_id,omitempty"`
	AckWaitingSeq     uint64 `json:"ack_waiting_sequence,omitempty"`
	RetriesRemaining  int32  `json:"retries_remaining"`
}

type stakingFile struct {
	BondDenom            string                `json:"bond_denom"`
	Delegations          []stakingRecord       `json:"delegations"`
	Unbondings           []stakingRecord       `json:"unbondings"`
	ShareclassShares     []shareclassRecord    `json:"shareclass_shares"`
	ShareclassUnbondings []shareclassUnbonding `json:"shareclass_unbondings"`
}

type stakingRecord struct {
	Owner     string `json:"owner"`
	Validator string `json:"validator"`
	Denom     string `json:"denom"`
	Amount    string `json:"amount"`
}

type shareclassRecord struct {
	Owner       string `json:"owner"`
	Validator   string `json:"validator"`
	ShareDenom  string `json:"share_denom"`
	ShareAmount string `json:"share_amount"`
	Denom       string `json:"denom"`
	Amount      string `json:"amount"`
}

type shareclassUnbonding struct {
	Delegator string `json:"delegator"`
	Recipient string `json:"recipient"`
	Validator string `json:"validator"`
	Denom     string `json:"denom"`
	Amount    string `json:"amount"`
}

type poolExport struct {
	Pools                []exportedPool                  `json:"pools"`
	Positions            []exportedPosition              `json:"positions"`
	AccumulatorPositions []pooltypes.AccumulatorPosition `json:"accumulator_positions"`
}

type exportedPool struct {
	ID               flexUint64           `json:"id"`
	DenomBase        string               `json:"denom_base"`
	DenomQuote       string               `json:"denom_quote"`
	TickParams       pooltypes.TickParams `json:"tick_params"`
	CurrentTick      flexInt64            `json:"current_tick"`
	CurrentSqrtPrice string               `json:"current_sqrt_price"`
}

type exportedPosition struct {
	ID        flexUint64 `json:"id"`
	Address   string     `json:"address"`
	PoolID    flexUint64 `json:"pool_id"`
	LowerTick flexInt64  `json:"lower_tick"`
	UpperTick flexInt64  `json:"upper_tick"`
	Liquidity string     `json:"liquidity"`
}

type flexUint64 uint64

func (n *flexUint64) UnmarshalJSON(data []byte) error {
	parsed, err := parseIntegerString(data)
	if err != nil {
		return fmt.Errorf("flexUint64: %w", err)
	}
	if parsed < 0 {
		return fmt.Errorf("flexUint64: negative value %d", parsed)
	}
	*n = flexUint64(parsed)
	return nil
}

type flexInt64 int64

func (n *flexInt64) UnmarshalJSON(data []byte) error {
	parsed, err := parseIntegerString(data)
	if err != nil {
		return fmt.Errorf("flexInt64: %w", err)
	}
	*n = flexInt64(parsed)
	return nil
}

func parseIntegerString(data []byte) (int64, error) {
	text := strings.Trim(string(data), `"`)
	if text == "null" || text == "" {
		return 0, nil
	}
	return strconv.ParseInt(text, 10, 64)
}

type lockupExport struct {
	LockupAccounts []exportedLockup `json:"lockup_accounts"`
}

type exportedLockup struct {
	Address string     `json:"address"`
	Owner   string     `json:"owner"`
	ID      flexUint64 `json:"id"`
}

type shareclassExport struct {
	Unbondings []exportedShareclassUnbonding `json:"unbondings"`
}

// exportedShareclassUnbonding accepts the upgrade export, which writes uint64 ids as strings.
type exportedShareclassUnbonding struct {
	ID               flexUint64 `json:"id"`
	RecipientAddress string     `json:"recipient_address"`
	DelegatorAddress string     `json:"delegator_address"`
	ValidatorAddress string     `json:"validator_address"`
	Amount           sdk.Coin   `json:"amount"`
}

type stakingExport struct {
	Params struct {
		BondDenom string `json:"bond_denom"`
	} `json:"params"`
	Validators []struct {
		OperatorAddress string         `json:"operator_address"`
		Tokens          math.Int       `json:"tokens"`
		DelegatorShares math.LegacyDec `json:"delegator_shares"`
	} `json:"validators"`
	Delegations []struct {
		DelegatorAddress string         `json:"delegator_address"`
		ValidatorAddress string         `json:"validator_address"`
		Shares           math.LegacyDec `json:"shares"`
	} `json:"delegations"`
	UnbondingDelegations []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Entries          []struct {
			Balance math.Int `json:"balance"`
		} `json:"entries"`
	} `json:"unbonding_delegations"`
}

func buildFiles(appState map[string]json.RawMessage) (bankFile, []positionRecord, []lockupRecord, inFlightFile, stakingFile, error) {
	var bankGenesis banktypes.GenesisState
	if err := unmarshalModule(appState, banktypes.ModuleName, &bankGenesis); err != nil {
		return bankFile{}, nil, nil, inFlightFile{}, stakingFile{}, err
	}
	var pools poolExport
	if err := unmarshalModule(appState, "liquiditypool", &pools); err != nil {
		return bankFile{}, nil, nil, inFlightFile{}, stakingFile{}, err
	}
	var lockup lockupExport
	if err := unmarshalModule(appState, lockuptypes.ModuleName, &lockup); err != nil {
		return bankFile{}, nil, nil, inFlightFile{}, stakingFile{}, err
	}
	var swap swapExport
	if err := unmarshalModule(appState, swaptypes.ModuleName, &swap); err != nil {
		return bankFile{}, nil, nil, inFlightFile{}, stakingFile{}, err
	}
	var shareclass shareclassExport
	if err := unmarshalModule(appState, shareclasstypes.ModuleName, &shareclass); err != nil {
		return bankFile{}, nil, nil, inFlightFile{}, stakingFile{}, err
	}
	var staking stakingExport
	if raw, ok := appState["staking"]; ok && len(bytes.TrimSpace(raw)) > 0 {
		if err := json.Unmarshal(raw, &staking); err != nil {
			return bankFile{}, nil, nil, inFlightFile{}, stakingFile{}, fmt.Errorf("buildFiles: unmarshal staking: %w", err)
		}
	}

	balancesByAddress := map[string]sdk.Coins{}
	for _, balance := range bankGenesis.Balances {
		balancesByAddress[balance.Address] = balance.Coins
	}

	positions, err := buildPositions(pools)
	if err != nil {
		return bankFile{}, nil, nil, inFlightFile{}, stakingFile{}, err
	}
	lockups := buildLockups(lockup, balancesByAddress)
	inFlight := buildInFlight(swap)
	stakingRecords, err := buildStaking(staking, bankGenesis, shareclass)
	if err != nil {
		return bankFile{}, nil, nil, inFlightFile{}, stakingFile{}, err
	}

	sort.Slice(bankGenesis.Balances, func(i, j int) bool {
		return bankGenesis.Balances[i].Address < bankGenesis.Balances[j].Address
	})
	return bankFile{Balances: bankGenesis.Balances}, positions, lockups, inFlight, stakingRecords, nil
}

func unmarshalModule(appState map[string]json.RawMessage, name string, dest any) error {
	raw, ok := appState[name]
	if !ok || len(bytes.TrimSpace(raw)) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, dest); err != nil {
		return fmt.Errorf("buildFiles: unmarshal %s: %w", name, err)
	}
	return nil
}

func buildPositions(genesis poolExport) ([]positionRecord, error) {
	poolsByID := make(map[uint64]exportedPool, len(genesis.Pools))
	for _, pool := range genesis.Pools {
		poolsByID[uint64(pool.ID)] = pool
	}
	rewardsByPosition := map[uint64]sdk.DecCoins{}
	for _, item := range genesis.AccumulatorPositions {
		var positionID uint64
		if _, err := fmt.Sscanf(item.Index, "fee_position_accumulator/|%d", &positionID); err != nil {
			continue
		}
		rewardsByPosition[positionID] = item.UnclaimedRewardsTotal
	}

	records := make([]positionRecord, 0, len(genesis.Positions))
	for _, position := range genesis.Positions {
		pool, found := poolsByID[uint64(position.PoolID)]
		if !found {
			return nil, fmt.Errorf("buildPositions: position %d references missing pool %d", position.ID, position.PoolID)
		}
		liquidity, err := math.LegacyNewDecFromStr(position.Liquidity)
		if err != nil {
			return nil, fmt.Errorf("buildPositions: position %d liquidity %q: %w", position.ID, position.Liquidity, err)
		}
		amountBase, amountQuote, err := withdrawalAmounts(toPool(pool), int64(position.LowerTick), int64(position.UpperTick), liquidity)
		if err != nil {
			return nil, fmt.Errorf("buildPositions: position %d pool %d: %w", position.ID, position.PoolID, err)
		}
		rewards := rewardsByPosition[uint64(position.ID)]
		if rewards == nil {
			rewards = sdk.DecCoins{}
		}
		records = append(records, positionRecord{
			PositionID:       uint64(position.ID),
			Owner:            position.Address,
			PoolID:           uint64(position.PoolID),
			LowerTick:        int64(position.LowerTick),
			UpperTick:        int64(position.UpperTick),
			Liquidity:        position.Liquidity,
			DenomBase:        pool.DenomBase,
			DenomQuote:       pool.DenomQuote,
			AmountBase:       amountBase.String(),
			AmountQuote:      amountQuote.String(),
			UnclaimedRewards: rewards,
		})
	}
	sort.Slice(records, func(i, j int) bool { return records[i].PositionID < records[j].PositionID })
	return records, nil
}

func toPool(pool exportedPool) pooltypes.Pool {
	return pooltypes.Pool{
		Id:               uint64(pool.ID),
		DenomBase:        pool.DenomBase,
		DenomQuote:       pool.DenomQuote,
		TickParams:       pool.TickParams,
		CurrentTick:      int64(pool.CurrentTick),
		CurrentSqrtPrice: pool.CurrentSqrtPrice,
	}
}

// withdrawalAmounts uses the same rounding as DecreaseLiquidity.
// Passing a negative liquidity delta rounds toward the pool, then the absolute value is the withdrawal amount.
func withdrawalAmounts(pool pooltypes.Pool, lowerTick, upperTick int64, liquidity math.LegacyDec) (math.Int, math.Int, error) {
	if liquidity.IsZero() {
		return math.ZeroInt(), math.ZeroInt(), nil
	}
	amountBase, amountQuote, err := pool.CalcActualAmounts(lowerTick, upperTick, liquidity.Neg())
	if err != nil {
		return math.Int{}, math.Int{}, err
	}
	return amountBase.TruncateInt().Abs(), amountQuote.TruncateInt().Abs(), nil
}

func buildLockups(genesis lockupExport, balancesByAddress map[string]sdk.Coins) []lockupRecord {
	records := make([]lockupRecord, 0, len(genesis.LockupAccounts))
	for _, account := range genesis.LockupAccounts {
		balances := balancesByAddress[account.Address]
		if balances == nil {
			balances = sdk.Coins{}
		}
		records = append(records, lockupRecord{
			Owner:         account.Owner,
			LockupAddress: account.Address,
			ID:            uint64(account.ID),
			Balances:      balances,
		})
	}
	sort.Slice(records, func(i, j int) bool {
		if records[i].Owner == records[j].Owner {
			return records[i].ID < records[j].ID
		}
		return records[i].Owner < records[j].Owner
	})
	return records
}

func buildInFlight(genesis swapExport) inFlightFile {
	incoming := make([]incomingPacketRecord, 0, len(genesis.Incoming))
	for _, packet := range genesis.Incoming {
		record := incomingPacketRecord{
			PortID:    packet.Index.PortID,
			ChannelID: packet.Index.ChannelID,
			Sequence:  uint64(packet.Index.Sequence),
		}
		var data transfertypes.FungibleTokenPacketData
		if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data); err == nil {
			record.Denom = data.Denom
			record.Amount = data.Amount
			record.Sender = data.Sender
			record.Receiver = data.Receiver
		}
		incoming = append(incoming, record)
	}
	outgoing := make([]outgoingPacketRecord, 0, len(genesis.Outgoing))
	for _, packet := range genesis.Outgoing {
		outgoing = append(outgoing, outgoingPacketRecord{
			PortID:            packet.Index.PortID,
			ChannelID:         packet.Index.ChannelID,
			Sequence:          uint64(packet.Index.Sequence),
			AckWaitingPortID:  packet.AckWaitingIndex.PortID,
			AckWaitingChannel: packet.AckWaitingIndex.ChannelID,
			AckWaitingSeq:     uint64(packet.AckWaitingIndex.Sequence),
			RetriesRemaining:  int32(packet.RetriesRemaining),
		})
	}
	return inFlightFile{Incoming: incoming, Outgoing: outgoing}
}

func buildStaking(staking stakingExport, bank banktypes.GenesisState, shareclass shareclassExport) (stakingFile, error) {
	moduleAddress := authtypes.NewModuleAddress(shareclasstypes.ModuleName).String()
	validators := map[string]struct {
		tokens math.Int
		shares math.LegacyDec
	}{}
	for _, validator := range staking.Validators {
		validators[validator.OperatorAddress] = struct {
			tokens math.Int
			shares math.LegacyDec
		}{tokens: validator.Tokens, shares: validator.DelegatorShares}
	}

	file := stakingFile{BondDenom: staking.Params.BondDenom}
	for _, delegation := range staking.Delegations {
		if delegation.DelegatorAddress == moduleAddress {
			continue
		}
		amount, err := delegatedTokens(validators, delegation.ValidatorAddress, delegation.Shares)
		if err != nil {
			return stakingFile{}, fmt.Errorf("buildStaking: delegator %s validator %s: %w", delegation.DelegatorAddress, delegation.ValidatorAddress, err)
		}
		file.Delegations = append(file.Delegations, stakingRecord{
			Owner:     delegation.DelegatorAddress,
			Validator: delegation.ValidatorAddress,
			Denom:     staking.Params.BondDenom,
			Amount:    amount.String(),
		})
	}
	for _, unbonding := range staking.UnbondingDelegations {
		if unbonding.DelegatorAddress == moduleAddress {
			continue
		}
		total := math.ZeroInt()
		for _, entry := range unbonding.Entries {
			if entry.Balance.IsNil() {
				continue
			}
			total = total.Add(entry.Balance)
		}
		file.Unbondings = append(file.Unbondings, stakingRecord{
			Owner:     unbonding.DelegatorAddress,
			Validator: unbonding.ValidatorAddress,
			Denom:     staking.Params.BondDenom,
			Amount:    total.String(),
		})
	}

	supply := map[string]math.Int{}
	for _, coin := range bank.Supply {
		supply[coin.Denom] = coin.Amount
	}
	moduleStake := map[string]math.Int{}
	for _, delegation := range staking.Delegations {
		if delegation.DelegatorAddress != moduleAddress {
			continue
		}
		amount, err := delegatedTokens(validators, delegation.ValidatorAddress, delegation.Shares)
		if err != nil {
			return stakingFile{}, fmt.Errorf("buildStaking: shareclass module delegation to %s: %w", delegation.ValidatorAddress, err)
		}
		moduleStake[delegation.ValidatorAddress] = amount
	}

	sharePattern := shareclasstypes.NonVotingShareTokenDenomRegexp()
	for _, balance := range bank.Balances {
		for _, coin := range balance.Coins {
			matches := sharePattern.FindStringSubmatch(coin.Denom)
			if matches == nil {
				continue
			}
			validator := matches[1]
			totalShare := supply[coin.Denom]
			if totalShare.IsNil() || totalShare.IsZero() {
				return stakingFile{}, fmt.Errorf("buildStaking: supply for %s is zero while %s holds %s", coin.Denom, balance.Address, coin.Amount)
			}
			totalStaked, found := moduleStake[validator]
			if !found {
				return stakingFile{}, fmt.Errorf("buildStaking: shareclass module has no delegation to %s for denom %s", validator, coin.Denom)
			}
			amount, err := shareclasstypes.CalculateAmountByShare(totalShare, totalStaked, coin.Amount)
			if err != nil {
				return stakingFile{}, fmt.Errorf("buildStaking: share %s of %s: %w", coin.Amount, coin.Denom, err)
			}
			file.ShareclassShares = append(file.ShareclassShares, shareclassRecord{
				Owner:       balance.Address,
				Validator:   validator,
				ShareDenom:  coin.Denom,
				ShareAmount: coin.Amount.String(),
				Denom:       staking.Params.BondDenom,
				Amount:      amount.String(),
			})
		}
	}

	for _, unbonding := range shareclass.Unbondings {
		file.ShareclassUnbondings = append(file.ShareclassUnbondings, shareclassUnbonding{
			Delegator: unbonding.DelegatorAddress,
			Recipient: unbonding.RecipientAddress,
			Validator: unbonding.ValidatorAddress,
			Denom:     unbonding.Amount.Denom,
			Amount:    unbonding.Amount.Amount.String(),
		})
	}
	if file.Delegations == nil {
		file.Delegations = []stakingRecord{}
	}
	if file.Unbondings == nil {
		file.Unbondings = []stakingRecord{}
	}
	if file.ShareclassShares == nil {
		file.ShareclassShares = []shareclassRecord{}
	}
	if file.ShareclassUnbondings == nil {
		file.ShareclassUnbondings = []shareclassUnbonding{}
	}
	return file, nil
}

func delegatedTokens(validators map[string]struct {
	tokens math.Int
	shares math.LegacyDec
}, validator string, shares math.LegacyDec) (math.Int, error) {
	info, found := validators[validator]
	if !found {
		return math.Int{}, fmt.Errorf("validator %s is not in the export", validator)
	}
	if info.shares.IsNil() || info.shares.IsZero() {
		return math.ZeroInt(), nil
	}
	if info.tokens.IsNil() {
		return math.Int{}, fmt.Errorf("validator %s has no token amount", validator)
	}
	return shares.MulInt(info.tokens).Quo(info.shares).TruncateInt(), nil
}
