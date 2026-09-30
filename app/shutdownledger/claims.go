// claims.go builds claims.json from a pre-upgrade export.
//
// Each row is one payable amount for one sunrise account, asset, and source.
// The file is the input to a later Edge claim for RISE and to the external
// payout checker for USDC and Cosmos assets. It does not move tokens.
//
// The export must be the pre-upgrade snapshot taken before IBC balances are
// swept. Pool, pool-fee, transfer-escrow, and protocol module balances are
// omitted because those coins are represented by positions, in-flight packets,
// delegations, or reward formulas. The USDN held by the USDrise wrapper is
// omitted because uusdrise is paid as USDC instead.
package shutdownledger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	crontypes "github.com/sunriselayer/sunrise/x/cron/types"
	datypes "github.com/sunriselayer/sunrise/x/da/types"
	feetypes "github.com/sunriselayer/sunrise/x/fee/types"
	incentivetypes "github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	poolkeeper "github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	pooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	lockuptypes "github.com/sunriselayer/sunrise/x/lockup/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
	stabletypes "github.com/sunriselayer/sunrise/x/stable/types"
	swaptypes "github.com/sunriselayer/sunrise/x/swap/types"
	convertertypes "github.com/sunriselayer/sunrise/x/tokenconverter/types"
	factorytypes "github.com/sunriselayer/sunrise/x/tokenfactory/types"
)

const (
	// AssetRise is urise and uvrise, paid later as an Edge coin.
	AssetRise = "rise"
	// AssetUSDrise is uusdrise, paid as USDC.
	AssetUSDrise = "usdrise"
	// AssetUSDN is unwrapped IBC USDN, paid as USDC.
	AssetUSDN = "usdn"

	PayoutEdge   = "edge"
	PayoutUSDC   = "usdc"
	PayoutCosmos = "cosmos"

	SourceBank           = "bank"
	SourceDelegation     = "delegation"
	SourceUnbonding      = "unbonding"
	SourceShareclass     = "shareclass"
	SourceStakingReward  = "staking_reward"
	SourcePosition       = "position"
	SourcePositionReward = "position_reward"
	SourceLockup         = "lockup"

	denomURise    = "urise"
	denomUVRise   = "uvrise"
	denomUUSDRise = "uusdrise"

	// usdnIBCDenom is the transfer denom of USDN on sunrise-1.
	usdnIBCDenom = "ibc/A7AD825A4B48DDA0138D118655E60100D22A4D690C45B95221520B58C9A64B63"
	// usdriseWrapper is the USDRise contract that custodies wrapped USDN.
	// Its USDN balance is the backing for uusdrise and is not a second claim.
	usdriseWrapper = "sunrise14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s2v9j75"

	slippageNote = "usdrise and usdn amounts are ledger base units. USDC sent from the hot wallet can be lower after the USDN to USDC swap."
)

// Claim is one payable row in claims.json.
type Claim struct {
	Owner       string `json:"owner"`
	Asset       string `json:"asset"`
	Amount      string `json:"amount"`
	Source      string `json:"source"`
	SourceID    string `json:"source_id"`
	ClaimableAt int64  `json:"claimable_at"`
	Payout      string `json:"payout"`
}

type claimsFile struct {
	SnapshotTime   string  `json:"snapshot_time"`
	SnapshotUnix   int64   `json:"snapshot_unix"`
	SnapshotHeight int64   `json:"snapshot_height"`
	Slippage       string  `json:"slippage"`
	Claims         []Claim `json:"claims"`
}

type exportDocument struct {
	GenesisTime   string                     `json:"genesis_time"`
	InitialHeight json.RawMessage            `json:"initial_height"`
	AppState      map[string]json.RawMessage `json:"app_state"`
}

type claimsPoolExport struct {
	Pools                []exportedPool                  `json:"pools"`
	Positions            []exportedPosition              `json:"positions"`
	AccumulatorPositions []pooltypes.AccumulatorPosition `json:"accumulator_positions"`
	Accumulators         []pooltypes.AccumulatorObject   `json:"accumulators"`
}

type claimsLockupExport struct {
	LockupAccounts []exportedLockupAccount `json:"lockup_accounts"`
}

type exportedLockupAccount struct {
	Address           string     `json:"address"`
	Owner             string     `json:"owner"`
	ID                flexUint64 `json:"id"`
	StartTime         flexInt64  `json:"start_time"`
	EndTime           flexInt64  `json:"end_time"`
	OriginalLocking   math.Int   `json:"original_locking"`
	AdditionalLocking math.Int   `json:"additional_locking"`
	DelegatedFree     math.Int   `json:"delegated_free"`
	DelegatedLocking  math.Int   `json:"delegated_locking"`
}

type shareclassRewardExport struct {
	Unbondings  []shareclasstypes.Unbonding `json:"unbondings"`
	Multipliers []struct {
		Validator        string `json:"validator"`
		Denom            string `json:"denom"`
		RewardMultiplier string `json:"reward_multiplier"`
	} `json:"reward_multipliers"`
	UserLast []struct {
		User             string `json:"user"`
		Validator        string `json:"validator"`
		Denom            string `json:"denom"`
		RewardMultiplier string `json:"reward_multiplier"`
	} `json:"user_last_reward_multipliers"`
}

type claimBuilder struct {
	snapshot     int64
	lockups      map[string]lockuptypes.LockupAccount
	excluded     map[string]struct{}
	rewards      *rewardState
	globalMult   map[string]map[string]math.LegacyDec
	userLastMult map[string]math.LegacyDec
	virtualMult  map[string]map[string]math.LegacyDec
	claims       []Claim
}

func buildClaims(doc exportDocument) (claimsFile, error) {
	snapshot, height, err := snapshotMeta(doc)
	if err != nil {
		return claimsFile{}, err
	}

	var bank banktypes.GenesisState
	if err := unmarshalModule(doc.AppState, banktypes.ModuleName, &bank); err != nil {
		return claimsFile{}, err
	}
	var pools claimsPoolExport
	if err := unmarshalModule(doc.AppState, pooltypes.ModuleName, &pools); err != nil {
		return claimsFile{}, err
	}
	var lockup claimsLockupExport
	if err := unmarshalModule(doc.AppState, lockuptypes.ModuleName, &lockup); err != nil {
		return claimsFile{}, err
	}
	var staking stakingExport
	if raw, ok := doc.AppState["staking"]; ok && len(bytes.TrimSpace(raw)) > 0 {
		if err := json.Unmarshal(raw, &staking); err != nil {
			return claimsFile{}, fmt.Errorf("buildClaims: unmarshal staking: %w", err)
		}
	}
	var shareclass shareclassRewardExport
	if err := unmarshalModule(doc.AppState, shareclasstypes.ModuleName, &shareclass); err != nil {
		return claimsFile{}, err
	}
	var distribution distributionExport
	if err := unmarshalModule(doc.AppState, disttypes.ModuleName, &distribution); err != nil {
		return claimsFile{}, err
	}
	rewards, err := prepareRewards(distribution, staking, height)
	if err != nil {
		return claimsFile{}, err
	}
	excluded, err := excludedAddresses(pools, staking, doc.AppState)
	if err != nil {
		return claimsFile{}, err
	}
	globalMult, err := multiplierIndex(shareclass)
	if err != nil {
		return claimsFile{}, err
	}
	userLast, err := userLastIndex(shareclass)
	if err != nil {
		return claimsFile{}, err
	}

	builder := &claimBuilder{
		snapshot:     snapshot,
		lockups:      lockupIndex(lockup),
		excluded:     excluded,
		rewards:      rewards,
		globalMult:   globalMult,
		userLastMult: userLast,
		virtualMult:  map[string]map[string]math.LegacyDec{},
		claims:       []Claim{},
	}
	if err := builder.addBank(bank); err != nil {
		return claimsFile{}, err
	}
	if err := builder.addPositions(pools); err != nil {
		return claimsFile{}, err
	}
	if err := builder.addStaking(staking, bank, shareclass); err != nil {
		return claimsFile{}, err
	}
	if err := builder.addShareclass(staking, bank, shareclass); err != nil {
		return claimsFile{}, err
	}
	builder.sortClaims()
	return claimsFile{
		SnapshotTime:   doc.GenesisTime,
		SnapshotUnix:   snapshot,
		SnapshotHeight: height,
		Slippage:       slippageNote,
		Claims:         builder.claims,
	}, nil
}

func snapshotMeta(doc exportDocument) (int64, int64, error) {
	var snapshot int64
	if doc.GenesisTime != "" {
		parsed, err := time.Parse(time.RFC3339Nano, doc.GenesisTime)
		if err != nil {
			return 0, 0, fmt.Errorf("snapshotMeta: genesis_time %q: %w", doc.GenesisTime, err)
		}
		snapshot = parsed.Unix()
	}
	if len(bytes.TrimSpace(doc.InitialHeight)) == 0 || string(doc.InitialHeight) == "null" {
		return snapshot, 0, nil
	}
	height, err := parseIntegerString(doc.InitialHeight)
	if err != nil {
		return 0, 0, fmt.Errorf("snapshotMeta: initial_height %s: %w", doc.InitialHeight, err)
	}
	return snapshot, height, nil
}

func lockupIndex(genesis claimsLockupExport) map[string]lockuptypes.LockupAccount {
	index := make(map[string]lockuptypes.LockupAccount, len(genesis.LockupAccounts))
	for _, account := range genesis.LockupAccounts {
		index[account.Address] = lockuptypes.LockupAccount{
			Address:           account.Address,
			Owner:             account.Owner,
			Id:                uint64(account.ID),
			StartTime:         int64(account.StartTime),
			EndTime:           int64(account.EndTime),
			OriginalLocking:   nonNilInt(account.OriginalLocking),
			AdditionalLocking: nonNilInt(account.AdditionalLocking),
			DelegatedFree:     nonNilInt(account.DelegatedFree),
			DelegatedLocking:  nonNilInt(account.DelegatedLocking),
		}
	}
	return index
}

func nonNilInt(value math.Int) math.Int {
	if value.IsNil() {
		return math.ZeroInt()
	}
	return value
}

func excludedAddresses(pools claimsPoolExport, staking stakingExport, appState map[string]json.RawMessage) (map[string]struct{}, error) {
	excluded := map[string]struct{}{}
	for _, name := range []string{
		authtypes.FeeCollectorName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		disttypes.ModuleName,
		minttypes.ModuleName,
		govtypes.ModuleName,
		transfertypes.ModuleName,
		shareclasstypes.ModuleName,
		pooltypes.ModuleName,
		lockuptypes.ModuleName,
		feetypes.ModuleName,
		incentivetypes.ModuleName,
		stabletypes.ModuleName,
		convertertypes.ModuleName,
		factorytypes.ModuleName,
		crontypes.ModuleName,
		datypes.ModuleName,
		swaptypes.ModuleName,
	} {
		excluded[authtypes.NewModuleAddress(name).String()] = struct{}{}
	}
	for _, pool := range pools.Pools {
		id := uint64(pool.ID)
		excluded[pooltypes.NewPoolAddress(id).String()] = struct{}{}
		excluded[pooltypes.NewPoolFeesAddress(id).String()] = struct{}{}
	}
	for _, validator := range staking.Validators {
		excluded[shareclasstypes.RewardSaverAddress(validator.OperatorAddress).String()] = struct{}{}
	}
	var ibc struct {
		ChannelGenesis struct {
			Channels []struct {
				PortID    string `json:"port_id"`
				ChannelID string `json:"channel_id"`
			} `json:"channels"`
		} `json:"channel_genesis"`
	}
	if raw, ok := appState["ibc"]; ok && len(bytes.TrimSpace(raw)) > 0 {
		if err := json.Unmarshal(raw, &ibc); err != nil {
			return nil, fmt.Errorf("excludedAddresses: unmarshal ibc: %w", err)
		}
		for _, channel := range ibc.ChannelGenesis.Channels {
			if channel.PortID != transfertypes.PortID {
				continue
			}
			excluded[transfertypes.GetEscrowAddress(channel.PortID, channel.ChannelID).String()] = struct{}{}
		}
	}
	return excluded, nil
}

func multiplierIndex(genesis shareclassRewardExport) (map[string]map[string]math.LegacyDec, error) {
	index := map[string]map[string]math.LegacyDec{}
	for _, item := range genesis.Multipliers {
		parsed, err := parseDec(item.RewardMultiplier)
		if err != nil {
			return nil, fmt.Errorf("multiplierIndex: validator %s denom %s: %w", item.Validator, item.Denom, err)
		}
		if index[item.Validator] == nil {
			index[item.Validator] = map[string]math.LegacyDec{}
		}
		index[item.Validator][item.Denom] = parsed
	}
	return index, nil
}

func userLastIndex(genesis shareclassRewardExport) (map[string]math.LegacyDec, error) {
	index := map[string]math.LegacyDec{}
	for _, item := range genesis.UserLast {
		parsed, err := parseDec(item.RewardMultiplier)
		if err != nil {
			return nil, fmt.Errorf("userLastIndex: user %s validator %s denom %s: %w", item.User, item.Validator, item.Denom, err)
		}
		index[userLastKey(item.User, item.Validator, item.Denom)] = parsed
	}
	return index, nil
}

func userLastKey(user, validator, denom string) string {
	return user + "|" + validator + "|" + denom
}

func parseInt(value string) (math.Int, error) {
	parsed, ok := math.NewIntFromString(value)
	if !ok {
		return math.Int{}, fmt.Errorf("parseInt: %q", value)
	}
	return parsed, nil
}

func parseDec(value string) (math.LegacyDec, error) {
	if value == "" {
		return math.LegacyZeroDec(), nil
	}
	parsed, err := math.LegacyNewDecFromStr(value)
	if err != nil {
		return math.LegacyDec{}, fmt.Errorf("parseDec: %q: %w", value, err)
	}
	return parsed, nil
}

func (b *claimBuilder) add(owner, asset string, amount math.Int, source, sourceID, payout string, claimableAt int64) error {
	if amount.IsNil() || !amount.IsPositive() {
		return nil
	}
	if owner == "" {
		return fmt.Errorf("add: empty owner for %s %s %s", asset, source, sourceID)
	}
	b.claims = append(b.claims, Claim{
		Owner:       owner,
		Asset:       asset,
		Amount:      amount.String(),
		Source:      source,
		SourceID:    sourceID,
		ClaimableAt: claimableAt,
		Payout:      payout,
	})
	return nil
}

func classifyDenom(denom string) (asset, payout string, include bool) {
	switch denom {
	case denomURise, denomUVRise:
		return AssetRise, PayoutEdge, true
	case denomUUSDRise:
		return AssetUSDrise, PayoutUSDC, true
	case usdnIBCDenom:
		return AssetUSDN, PayoutUSDC, true
	default:
		if strings.HasPrefix(denom, "shareclass/") {
			return "", "", false
		}
		if strings.HasPrefix(denom, "ibc/") {
			return denom, PayoutCosmos, true
		}
		return denom, PayoutCosmos, true
	}
}

func (b *claimBuilder) addBank(bank banktypes.GenesisState) error {
	for _, balance := range bank.Balances {
		if _, skip := b.excluded[balance.Address]; skip {
			continue
		}
		for _, coin := range balance.Coins {
			if balance.Address == usdriseWrapper && coin.Denom == usdnIBCDenom {
				continue
			}
			asset, payout, include := classifyDenom(coin.Denom)
			if !include {
				continue
			}
			account, lockedAccount := b.lockups[balance.Address]
			if lockedAccount && asset == AssetRise {
				if err := b.addLockupLiquid(account, coin.Amount); err != nil {
					return err
				}
				continue
			}
			owner := balance.Address
			source := SourceBank
			sourceID := coin.Denom
			if lockedAccount {
				owner = account.Owner
				source = SourceLockup
				sourceID = fmt.Sprintf("%d:bank:%s", account.Id, coin.Denom)
			}
			if err := b.add(owner, asset, coin.Amount, source, sourceID, payout, b.snapshot); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *claimBuilder) addLockupLiquid(account lockuptypes.LockupAccount, amount math.Int) error {
	sourceID := fmt.Sprintf("%d:bank:%s", account.Id, denomURise)
	if !vestingOpen(account, b.snapshot) {
		return b.add(account.Owner, AssetRise, amount, SourceLockup, sourceID, PayoutEdge, b.snapshot)
	}
	_, locked, err := lockAmounts(account, b.snapshot)
	if err != nil {
		return fmt.Errorf("addLockupLiquid: lockup %d owner %s: %w", account.Id, account.Owner, err)
	}
	lockedLiquid := account.GetNotBondedLockedAmount(locked)
	if lockedLiquid.GT(amount) {
		lockedLiquid = amount
	}
	unlockedLiquid := amount.Sub(lockedLiquid)
	if err := b.add(account.Owner, AssetRise, lockedLiquid, SourceLockup, sourceID+":locked", PayoutEdge, account.EndTime); err != nil {
		return err
	}
	return b.add(account.Owner, AssetRise, unlockedLiquid, SourceLockup, sourceID+":unlocked", PayoutEdge, b.snapshot)
}

func vestingOpen(account lockuptypes.LockupAccount, snapshot int64) bool {
	return snapshot < account.EndTime
}

func lockAmounts(account lockuptypes.LockupAccount, snapshot int64) (math.Int, math.Int, error) {
	if account.StartTime == account.EndTime {
		total := account.OriginalLocking.Add(account.AdditionalLocking)
		if snapshot >= account.EndTime {
			return total, math.ZeroInt(), nil
		}
		return math.ZeroInt(), total, nil
	}
	unlocked, locked, err := account.GetLockCoinInfo(snapshot)
	if err != nil {
		return math.Int{}, math.Int{}, fmt.Errorf("lockAmounts: lockup %d: %w", account.Id, err)
	}
	return nonNilInt(unlocked), nonNilInt(locked), nil
}

func (b *claimBuilder) addPositions(pools claimsPoolExport) error {
	accumulators := map[string]pooltypes.AccumulatorObject{}
	for _, accumulator := range pools.Accumulators {
		accumulators[accumulator.Name] = accumulator
	}
	positions, err := buildPositions(poolExport{
		Pools:                pools.Pools,
		Positions:            pools.Positions,
		AccumulatorPositions: pools.AccumulatorPositions,
	})
	if err != nil {
		return err
	}
	rewardByID := map[uint64]pooltypes.AccumulatorPosition{}
	for _, item := range pools.AccumulatorPositions {
		var positionID uint64
		if _, err := fmt.Sscanf(item.Index, "fee_position_accumulator/|%d", &positionID); err != nil {
			continue
		}
		rewardByID[positionID] = item
	}
	for _, position := range positions {
		owner := position.Owner
		if account, ok := b.lockups[owner]; ok {
			owner = account.Owner
		}
		base, err := parseInt(position.AmountBase)
		if err != nil {
			return fmt.Errorf("addPositions: position %d base %q: %w", position.PositionID, position.AmountBase, err)
		}
		quote, err := parseInt(position.AmountQuote)
		if err != nil {
			return fmt.Errorf("addPositions: position %d quote %q: %w", position.PositionID, position.AmountQuote, err)
		}
		if err := b.addClassified(owner, position.DenomBase, base, SourcePosition, fmt.Sprintf("%d:base", position.PositionID)); err != nil {
			return err
		}
		if err := b.addClassified(owner, position.DenomQuote, quote, SourcePosition, fmt.Sprintf("%d:quote", position.PositionID)); err != nil {
			return err
		}
		rewards := positionRewardCoins(accumulators, position.PoolID, rewardByID[position.PositionID])
		for _, coin := range rewards {
			sourceID := fmt.Sprintf("%d:%s", position.PositionID, coin.Denom)
			if err := b.addClassified(owner, coin.Denom, coin.Amount, SourcePositionReward, sourceID); err != nil {
				return err
			}
		}
	}
	return nil
}

func positionRewardCoins(accumulators map[string]pooltypes.AccumulatorObject, poolID uint64, position pooltypes.AccumulatorPosition) sdk.Coins {
	rewards := position.UnclaimedRewardsTotal
	if accumulator, found := accumulators[pooltypes.KeyFeePoolAccumulator(poolID)]; found && position.NumShares != "" {
		rewards = poolkeeper.GetTotalRewards(accumulator, position)
	}
	coins, _ := rewards.TruncateDecimal()
	return coins
}

func (b *claimBuilder) addClassified(owner, denom string, amount math.Int, source, sourceID string) error {
	asset, payout, include := classifyDenom(denom)
	if !include {
		return nil
	}
	return b.add(owner, asset, amount, source, sourceID, payout, b.snapshot)
}

func (b *claimBuilder) addStaking(staking stakingExport, bank banktypes.GenesisState, shareclass shareclassRewardExport) error {
	records, err := buildStaking(staking, bank, shareclassExport{Unbondings: shareclass.Unbondings})
	if err != nil {
		return err
	}
	grouped := map[string][]stakingRecord{}
	for _, record := range records.Delegations {
		if _, ok := b.lockups[record.Owner]; ok {
			grouped[record.Owner] = append(grouped[record.Owner], record)
			continue
		}
		amount, err := parseInt(record.Amount)
		if err != nil {
			return fmt.Errorf("addStaking: delegation %s to %s amount %q: %w", record.Owner, record.Validator, record.Amount, err)
		}
		if err := b.addClassified(record.Owner, record.Denom, amount, SourceDelegation, record.Validator); err != nil {
			return err
		}
	}
	for address, rows := range grouped {
		if err := b.addLockupDelegations(b.lockups[address], rows); err != nil {
			return err
		}
	}
	for _, record := range records.Unbondings {
		amount, err := parseInt(record.Amount)
		if err != nil {
			return fmt.Errorf("addStaking: unbonding %s from %s amount %q: %w", record.Owner, record.Validator, record.Amount, err)
		}
		if account, ok := b.lockups[record.Owner]; ok {
			if record.Denom != denomURise && record.Denom != denomUVRise {
				return fmt.Errorf("addStaking: lockup %d unbonding denom %s is not rise", account.Id, record.Denom)
			}
			if err := b.addLockupRise(account, amount, "unbonding:"+record.Validator, true); err != nil {
				return err
			}
			continue
		}
		if err := b.addClassified(record.Owner, record.Denom, amount, SourceUnbonding, record.Validator); err != nil {
			return err
		}
	}

	moduleAddress := authtypes.NewModuleAddress(shareclasstypes.ModuleName).String()
	for _, delegation := range staking.Delegations {
		if delegation.DelegatorAddress == moduleAddress {
			continue
		}
		coins, err := b.rewards.delegationRewards(delegation.DelegatorAddress, delegation.ValidatorAddress, delegation.Shares)
		if err != nil {
			return err
		}
		for _, coin := range coins {
			sourceID := delegation.ValidatorAddress + ":" + coin.Denom
			if account, ok := b.lockups[delegation.DelegatorAddress]; ok {
				if err := b.addMaybeLocked(account, coin.Denom, coin.Amount, "staking_reward:"+sourceID); err != nil {
					return err
				}
				continue
			}
			if err := b.addClassified(delegation.DelegatorAddress, coin.Denom, coin.Amount, SourceStakingReward, sourceID); err != nil {
				return err
			}
		}
	}
	for _, commission := range b.rewards.commissions {
		account, err := validatorAccount(commission.validator)
		if err != nil {
			return err
		}
		for _, coin := range commission.coins {
			sourceID := "commission:" + commission.validator + ":" + coin.Denom
			if err := b.addClassified(account, coin.Denom, coin.Amount, SourceStakingReward, sourceID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *claimBuilder) addLockupDelegations(account lockuptypes.LockupAccount, rows []stakingRecord) error {
	sort.Slice(rows, func(i, j int) bool { return rows[i].Validator < rows[j].Validator })
	remainingLocked := math.ZeroInt()
	if vestingOpen(account, b.snapshot) {
		remainingLocked = account.DelegatedLocking
	}
	for _, row := range rows {
		amount, err := parseInt(row.Amount)
		if err != nil {
			return fmt.Errorf("addLockupDelegations: lockup %d validator %s amount %q: %w", account.Id, row.Validator, row.Amount, err)
		}
		if row.Denom != denomURise && row.Denom != denomUVRise {
			return fmt.Errorf("addLockupDelegations: lockup %d validator %s denom %s is not rise", account.Id, row.Validator, row.Denom)
		}
		locked := math.ZeroInt()
		if remainingLocked.IsPositive() {
			locked = amount
			if locked.GT(remainingLocked) {
				locked = remainingLocked
			}
			remainingLocked = remainingLocked.Sub(locked)
		}
		unlocked := amount.Sub(locked)
		if err := b.add(account.Owner, AssetRise, locked, SourceLockup, fmt.Sprintf("%d:delegation:%s:locked", account.Id, row.Validator), PayoutEdge, account.EndTime); err != nil {
			return err
		}
		if err := b.add(account.Owner, AssetRise, unlocked, SourceLockup, fmt.Sprintf("%d:delegation:%s:unlocked", account.Id, row.Validator), PayoutEdge, b.snapshot); err != nil {
			return err
		}
	}
	return nil
}

func (b *claimBuilder) addLockupRise(account lockuptypes.LockupAccount, amount math.Int, sourceID string, followVest bool) error {
	fullID := fmt.Sprintf("%d:%s", account.Id, sourceID)
	if !followVest || !vestingOpen(account, b.snapshot) {
		return b.add(account.Owner, AssetRise, amount, SourceLockup, fullID, PayoutEdge, b.snapshot)
	}
	return b.add(account.Owner, AssetRise, amount, SourceLockup, fullID, PayoutEdge, account.EndTime)
}

func (b *claimBuilder) addMaybeLocked(account lockuptypes.LockupAccount, denom string, amount math.Int, sourceID string) error {
	asset, payout, include := classifyDenom(denom)
	if !include {
		return nil
	}
	nowAmount, laterAmount, laterAt, err := b.splitVest(account, amount)
	if err != nil {
		return fmt.Errorf("addMaybeLocked: lockup %d %s: %w", account.Id, sourceID, err)
	}
	fullID := fmt.Sprintf("%d:%s", account.Id, sourceID)
	if err := b.add(account.Owner, asset, nowAmount, SourceLockup, fullID+":unlocked", payout, b.snapshot); err != nil {
		return err
	}
	return b.add(account.Owner, asset, laterAmount, SourceLockup, fullID+":locked", payout, laterAt)
}

func (b *claimBuilder) splitVest(account lockuptypes.LockupAccount, amount math.Int) (math.Int, math.Int, int64, error) {
	if amount.IsNil() || !amount.IsPositive() || !vestingOpen(account, b.snapshot) {
		return amount, math.ZeroInt(), b.snapshot, nil
	}
	unlocked, locked, err := lockAmounts(account, b.snapshot)
	if err != nil {
		return math.Int{}, math.Int{}, 0, err
	}
	total := unlocked.Add(locked)
	if total.IsZero() {
		return amount, math.ZeroInt(), b.snapshot, nil
	}
	nowAmount := amount.Mul(unlocked).Quo(total)
	return nowAmount, amount.Sub(nowAmount), account.EndTime, nil
}

func (b *claimBuilder) addShareclass(staking stakingExport, bank banktypes.GenesisState, shareclass shareclassRewardExport) error {
	if err := b.foldModuleRewards(staking, bank); err != nil {
		return err
	}
	records, err := buildStaking(staking, bank, shareclassExport{Unbondings: shareclass.Unbondings})
	if err != nil {
		return err
	}
	for _, share := range records.ShareclassShares {
		underlying, err := parseInt(share.Amount)
		if err != nil {
			return fmt.Errorf("addShareclass: %s share %s: %w", share.Owner, share.ShareDenom, err)
		}
		shareAmount, err := parseInt(share.ShareAmount)
		if err != nil {
			return fmt.Errorf("addShareclass: %s share amount %s: %w", share.Owner, share.ShareAmount, err)
		}
		if account, ok := b.lockups[share.Owner]; ok {
			if err := b.addMaybeLocked(account, share.Denom, underlying, "shareclass:"+share.Validator); err != nil {
				return err
			}
			if err := b.addShareRewards(account.Owner, &account, share.Validator, shareAmount); err != nil {
				return err
			}
			continue
		}
		if err := b.addClassified(share.Owner, share.Denom, underlying, SourceShareclass, share.Validator); err != nil {
			return err
		}
		if err := b.addShareRewards(share.Owner, nil, share.Validator, shareAmount); err != nil {
			return err
		}
	}
	for _, unbonding := range shareclass.Unbondings {
		if unbonding.Amount.IsNil() || !unbonding.Amount.Amount.IsPositive() {
			continue
		}
		sourceID := unbonding.ValidatorAddress + ":unbonding"
		if account, ok := b.lockups[unbonding.RecipientAddress]; ok {
			if err := b.addMaybeLocked(account, unbonding.Amount.Denom, unbonding.Amount.Amount, "shareclass:"+sourceID); err != nil {
				return err
			}
			continue
		}
		if err := b.addClassified(unbonding.RecipientAddress, unbonding.Amount.Denom, unbonding.Amount.Amount, SourceShareclass, sourceID); err != nil {
			return err
		}
	}
	return nil
}

func (b *claimBuilder) foldModuleRewards(staking stakingExport, bank banktypes.GenesisState) error {
	moduleAddress := authtypes.NewModuleAddress(shareclasstypes.ModuleName).String()
	supply := map[string]math.Int{}
	for _, coin := range bank.Supply {
		supply[coin.Denom] = coin.Amount
	}
	for _, delegation := range staking.Delegations {
		if delegation.DelegatorAddress != moduleAddress {
			continue
		}
		coins, err := b.rewards.delegationRewards(delegation.DelegatorAddress, delegation.ValidatorAddress, delegation.Shares)
		if err != nil {
			return err
		}
		shareDenom := shareclasstypes.NonVotingShareTokenDenom(delegation.ValidatorAddress)
		totalShare := supply[shareDenom]
		if totalShare.IsNil() || !totalShare.IsPositive() {
			continue
		}
		totalDec, err := math.LegacyNewDecFromStr(totalShare.String())
		if err != nil {
			return fmt.Errorf("foldModuleRewards: supply %s: %w", shareDenom, err)
		}
		for _, coin := range coins {
			ratio, err := math.LegacyNewDecFromStr(coin.Amount.String())
			if err != nil {
				return fmt.Errorf("foldModuleRewards: reward %s: %w", coin, err)
			}
			current := b.effectiveMultiplier(delegation.ValidatorAddress, coin.Denom)
			if b.virtualMult[delegation.ValidatorAddress] == nil {
				b.virtualMult[delegation.ValidatorAddress] = map[string]math.LegacyDec{}
			}
			b.virtualMult[delegation.ValidatorAddress][coin.Denom] = current.Add(ratio.Quo(totalDec))
		}
	}
	return nil
}

func (b *claimBuilder) effectiveMultiplier(validator, denom string) math.LegacyDec {
	if denoms := b.virtualMult[validator]; denoms != nil {
		if value, ok := denoms[denom]; ok {
			return value
		}
	}
	if denoms := b.globalMult[validator]; denoms != nil {
		if value, ok := denoms[denom]; ok {
			return value
		}
	}
	return math.LegacyZeroDec()
}

func (b *claimBuilder) addShareRewards(owner string, lockup *lockuptypes.LockupAccount, validator string, shareAmount math.Int) error {
	denoms := map[string]struct{}{}
	for denom := range b.globalMult[validator] {
		denoms[denom] = struct{}{}
	}
	for denom := range b.virtualMult[validator] {
		denoms[denom] = struct{}{}
	}
	names := make([]string, 0, len(denoms))
	for denom := range denoms {
		names = append(names, denom)
	}
	sort.Strings(names)
	for _, denom := range names {
		last := b.userLastMult[userLastKey(owner, validator, denom)]
		if last.IsNil() {
			last = math.LegacyZeroDec()
		}
		holder := owner
		if lockup != nil {
			holder = lockup.Address
			last = b.userLastMult[userLastKey(lockup.Address, validator, denom)]
			if last.IsNil() {
				last = math.LegacyZeroDec()
			}
		}
		reward, err := shareclasstypes.CalculateReward(b.effectiveMultiplier(validator, denom), last, shareAmount)
		if err != nil {
			return fmt.Errorf("addShareRewards: %s validator %s denom %s share %s: %w", holder, validator, denom, shareAmount, err)
		}
		sourceID := "shareclass:" + validator + ":" + denom
		if lockup != nil {
			if err := b.addMaybeLocked(*lockup, denom, reward, sourceID); err != nil {
				return err
			}
			continue
		}
		if err := b.addClassified(owner, denom, reward, SourceStakingReward, sourceID); err != nil {
			return err
		}
	}
	return nil
}

func (b *claimBuilder) sortClaims() {
	sort.Slice(b.claims, func(i, j int) bool {
		if b.claims[i].Owner != b.claims[j].Owner {
			return b.claims[i].Owner < b.claims[j].Owner
		}
		if b.claims[i].Asset != b.claims[j].Asset {
			return b.claims[i].Asset < b.claims[j].Asset
		}
		if b.claims[i].Source != b.claims[j].Source {
			return b.claims[i].Source < b.claims[j].Source
		}
		if b.claims[i].SourceID != b.claims[j].SourceID {
			return b.claims[i].SourceID < b.claims[j].SourceID
		}
		if b.claims[i].ClaimableAt != b.claims[j].ClaimableAt {
			return b.claims[i].ClaimableAt < b.claims[j].ClaimableAt
		}
		return b.claims[i].Amount < b.claims[j].Amount
	})
}

func validatorAccount(valoper string) (string, error) {
	raw, err := sdk.ValAddressFromBech32(valoper)
	if err != nil {
		return "", fmt.Errorf("validatorAccount: %s: %w", valoper, err)
	}
	return sdk.AccAddress(raw).String(), nil
}
