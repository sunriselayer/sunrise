// Delegation rewards follow the Cosmos distribution period formula.
// The snapshot does not increment validator periods on chain. This file
// reproduces that increment from the exported current and historical rewards,
// then applies slash events between the delegation's starting height and the
// export height. Dust below one base unit is not claimed.
package shutdownledger

import (
	"fmt"
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type distributionExport struct {
	Historical []struct {
		ValidatorAddress string     `json:"validator_address"`
		Period           flexUint64 `json:"period"`
		Rewards          struct {
			Cumulative sdk.DecCoins `json:"cumulative_reward_ratio"`
		} `json:"rewards"`
	} `json:"validator_historical_rewards"`
	Current []struct {
		ValidatorAddress string `json:"validator_address"`
		Rewards          struct {
			Rewards sdk.DecCoins `json:"rewards"`
			Period  flexUint64   `json:"period"`
		} `json:"rewards"`
	} `json:"validator_current_rewards"`
	Starting []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		StartingInfo     struct {
			PreviousPeriod flexUint64     `json:"previous_period"`
			Stake          math.LegacyDec `json:"stake"`
			Height         flexUint64     `json:"creation_height"`
		} `json:"starting_info"`
	} `json:"delegator_starting_infos"`
	Slashes []struct {
		ValidatorAddress string     `json:"validator_address"`
		Height           flexUint64 `json:"height"`
		Event            struct {
			Period   flexUint64     `json:"validator_period"`
			Fraction math.LegacyDec `json:"fraction"`
		} `json:"validator_slash_event"`
	} `json:"validator_slash_events"`
	Commissions []struct {
		ValidatorAddress string `json:"validator_address"`
		Accumulated      struct {
			Commission sdk.DecCoins `json:"commission"`
		} `json:"accumulated"`
	} `json:"validator_accumulated_commissions"`
}

type startingInfo struct {
	period uint64
	stake  math.LegacyDec
	height uint64
}

type slashEvent struct {
	height   uint64
	period   uint64
	fraction math.LegacyDec
}

type currentReward struct {
	period     uint64
	cumulative sdk.DecCoins
}

type commissionReward struct {
	validator string
	coins     sdk.Coins
}

type rewardState struct {
	height     uint64
	validators map[string]struct {
		tokens math.Int
		shares math.LegacyDec
	}
	historical  map[string]map[uint64]sdk.DecCoins
	current     map[string]currentReward
	starting    map[string]startingInfo
	slashes     map[string][]slashEvent
	commissions []commissionReward
}

func prepareRewards(distribution distributionExport, staking stakingExport, height int64) (*rewardState, error) {
	if height < 0 {
		return nil, fmt.Errorf("prepareRewards: negative snapshot height %d", height)
	}
	state := &rewardState{
		height: uint64(height),
		validators: map[string]struct {
			tokens math.Int
			shares math.LegacyDec
		}{},
		historical: map[string]map[uint64]sdk.DecCoins{},
		current:    map[string]currentReward{},
		starting:   map[string]startingInfo{},
		slashes:    map[string][]slashEvent{},
	}
	for _, validator := range staking.Validators {
		state.validators[validator.OperatorAddress] = struct {
			tokens math.Int
			shares math.LegacyDec
		}{tokens: nonNilInt(validator.Tokens), shares: validator.DelegatorShares}
	}
	for _, item := range distribution.Historical {
		if state.historical[item.ValidatorAddress] == nil {
			state.historical[item.ValidatorAddress] = map[uint64]sdk.DecCoins{}
		}
		state.historical[item.ValidatorAddress][uint64(item.Period)] = item.Rewards.Cumulative
	}
	for _, item := range distribution.Starting {
		key := delegationKey(item.DelegatorAddress, item.ValidatorAddress)
		stake := item.StartingInfo.Stake
		if stake.IsNil() {
			stake = math.LegacyZeroDec()
		}
		state.starting[key] = startingInfo{
			period: uint64(item.StartingInfo.PreviousPeriod),
			stake:  stake,
			height: uint64(item.StartingInfo.Height),
		}
	}
	for _, item := range distribution.Slashes {
		state.slashes[item.ValidatorAddress] = append(state.slashes[item.ValidatorAddress], slashEvent{
			height:   uint64(item.Height),
			period:   uint64(item.Event.Period),
			fraction: item.Event.Fraction,
		})
	}
	for validator, events := range state.slashes {
		sort.SliceStable(events, func(i, j int) bool {
			return events[i].height < events[j].height
		})
		state.slashes[validator] = events
	}
	for _, item := range distribution.Current {
		period := uint64(item.Rewards.Period)
		if period == 0 {
			return nil, fmt.Errorf("prepareRewards: validator %s current period is 0", item.ValidatorAddress)
		}
		cumulative, err := state.closedCumulative(item.ValidatorAddress, period, item.Rewards.Rewards)
		if err != nil {
			return nil, err
		}
		state.current[item.ValidatorAddress] = currentReward{period: period, cumulative: cumulative}
	}
	for _, item := range distribution.Commissions {
		coins, _ := item.Accumulated.Commission.TruncateDecimal()
		if coins.IsZero() {
			continue
		}
		state.commissions = append(state.commissions, commissionReward{
			validator: item.ValidatorAddress,
			coins:     coins,
		})
	}
	return state, nil
}

func delegationKey(delegator, validator string) string {
	return delegator + "|" + validator
}

// closedCumulative is the historical ratio that IncrementValidatorPeriod would
// store for the current period. A zero-token validator does not add its current
// rewards; the protocol sends those to the community pool.
func (s *rewardState) closedCumulative(validator string, period uint64, current sdk.DecCoins) (sdk.DecCoins, error) {
	previous, err := s.storedCumulative(validator, period-1)
	if err != nil {
		return nil, fmt.Errorf("closedCumulative: validator %s: %w", validator, err)
	}
	info, found := s.validators[validator]
	if !found {
		return nil, fmt.Errorf("closedCumulative: validator %s is not in the staking export", validator)
	}
	if info.tokens.IsZero() {
		return previous, nil
	}
	tokenDec, err := math.LegacyNewDecFromStr(info.tokens.String())
	if err != nil {
		return nil, fmt.Errorf("closedCumulative: validator %s tokens %s: %w", validator, info.tokens, err)
	}
	if current == nil {
		current = sdk.DecCoins{}
	}
	ratio := current.QuoDecTruncate(tokenDec)
	return previous.Add(ratio...), nil
}

func (s *rewardState) storedCumulative(validator string, period uint64) (sdk.DecCoins, error) {
	periods := s.historical[validator]
	if periods == nil {
		return nil, fmt.Errorf("storedCumulative: validator %s has no historical rewards", validator)
	}
	coins, found := periods[period]
	if !found {
		return nil, fmt.Errorf("storedCumulative: validator %s period %d is missing", validator, period)
	}
	if coins == nil {
		return sdk.DecCoins{}, nil
	}
	return coins, nil
}

func (s *rewardState) cumulative(validator string, period uint64) (sdk.DecCoins, error) {
	if current, found := s.current[validator]; found && period == current.period {
		return current.cumulative, nil
	}
	return s.storedCumulative(validator, period)
}

func (s *rewardState) delegationRewards(delegator, validator string, shares math.LegacyDec) (sdk.Coins, error) {
	info, found := s.starting[delegationKey(delegator, validator)]
	if !found {
		return nil, fmt.Errorf("delegationRewards: delegator %s validator %s has no starting info", delegator, validator)
	}
	if info.height == s.height {
		return sdk.NewCoins(), nil
	}
	if info.height > s.height {
		return nil, fmt.Errorf("delegationRewards: delegator %s validator %s starts at height %d after snapshot %d", delegator, validator, info.height, s.height)
	}
	current, found := s.current[validator]
	if !found {
		return nil, fmt.Errorf("delegationRewards: validator %s has no current rewards", validator)
	}
	stake := info.stake
	startingPeriod := info.period
	var rewards sdk.DecCoins
	if s.height > info.height {
		for _, event := range s.slashes[validator] {
			if event.height < info.height || event.height > s.height {
				continue
			}
			if event.fraction.IsNil() {
				return nil, fmt.Errorf("delegationRewards: validator %s slash at height %d has no fraction", validator, event.height)
			}
			if event.period > startingPeriod {
				segment, err := s.rewardsBetween(validator, startingPeriod, event.period, stake)
				if err != nil {
					return nil, fmt.Errorf("delegationRewards: delegator %s validator %s: %w", delegator, validator, err)
				}
				rewards = rewards.Add(segment...)
				stake = stake.MulTruncate(math.LegacyOneDec().Sub(event.fraction))
				startingPeriod = event.period
			}
		}
	}
	validatorInfo, found := s.validators[validator]
	if !found {
		return nil, fmt.Errorf("delegationRewards: validator %s is not in the staking export", validator)
	}
	if shares.IsNil() {
		shares = math.LegacyZeroDec()
	}
	if validatorInfo.shares.IsNil() || validatorInfo.shares.IsZero() {
		return nil, fmt.Errorf("delegationRewards: validator %s has no delegator shares", validator)
	}
	currentStake := shares.MulInt(validatorInfo.tokens).Quo(validatorInfo.shares)
	if stake.GT(currentStake) {
		margin := math.LegacySmallestDec().MulInt64(3)
		if stake.LTE(currentStake.Add(margin)) {
			stake = currentStake
		} else {
			return nil, fmt.Errorf("delegationRewards: delegator %s validator %s calculated stake %s exceeds current stake %s", delegator, validator, stake, currentStake)
		}
	}
	finalSegment, err := s.rewardsBetween(validator, startingPeriod, current.period, stake)
	if err != nil {
		return nil, fmt.Errorf("delegationRewards: delegator %s validator %s: %w", delegator, validator, err)
	}
	rewards = rewards.Add(finalSegment...)
	coins, _ := rewards.TruncateDecimal()
	if coins == nil {
		return sdk.NewCoins(), nil
	}
	return coins, nil
}

func (s *rewardState) rewardsBetween(validator string, startingPeriod, endingPeriod uint64, stake math.LegacyDec) (sdk.DecCoins, error) {
	if startingPeriod > endingPeriod {
		return nil, fmt.Errorf("rewardsBetween: validator %s starting period %d is after ending period %d", validator, startingPeriod, endingPeriod)
	}
	if stake.IsNegative() {
		return nil, fmt.Errorf("rewardsBetween: validator %s stake %s is negative", validator, stake)
	}
	starting, err := s.cumulative(validator, startingPeriod)
	if err != nil {
		return nil, err
	}
	ending, err := s.cumulative(validator, endingPeriod)
	if err != nil {
		return nil, err
	}
	difference := ending.Sub(starting)
	if difference.IsAnyNegative() {
		return nil, fmt.Errorf("rewardsBetween: validator %s period %d to %d produced a negative ratio", validator, startingPeriod, endingPeriod)
	}
	return difference.MulDecTruncate(stake), nil
}
