package keeper

import (
	"context"
	"sort"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

// ValidatorGovInfo used for tallying
type ValidatorGovInfo struct {
	Address             sdk.ValAddress     // address of the validator operator
	BondedTokens        math.Int           // Power of a Validator
	DelegatorShares     math.LegacyDec     // Total outstanding delegator shares
	DelegatorDeductions math.LegacyDec     // Delegator deductions from validator's delegators voting independently
	PoolWeights         []types.PoolWeight // Vote of the validator
}

func (k Keeper) Tally(ctx context.Context) ([]types.TallyResult, error) {
	validators := make(map[string]ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	err := k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		valBz, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(validator.GetOperator())
		if err != nil {
			return false
		}
		validators[validator.GetOperator()] = ValidatorGovInfo{
			Address:             valBz,
			BondedTokens:        validator.GetBondedTokens(),
			DelegatorShares:     validator.GetDelegatorShares(),
			DelegatorDeductions: math.LegacyZeroDec(),
			PoolWeights:         []types.PoolWeight{},
		}

		return false
	})
	if err != nil {
		return []types.TallyResult{}, err
	}

	// Hereafter it is analogy with gov CalculateVoteResultsAndVotingPowerFn
	totalVotingPower := math.LegacyZeroDec()

	results := make(map[uint64]math.LegacyDec)

	// <sunrise>
	// Deduct shareclass module's delegations
	shareclassAddr := k.accountKeeper.GetModuleAddress(shareclasstypes.ModuleName)
	err = k.stakingKeeper.IterateDelegations(ctx, shareclassAddr, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
		valAddrStr := delegation.GetValidatorAddr()
		if val, ok := validators[valAddrStr]; ok {
			val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
			validators[valAddrStr] = val
		}
		return false
	})
	if err != nil {
		return []types.TallyResult{}, err
	}
	// </sunrise>

	votesToRemove := []sdk.AccAddress{}
	// iterate over all votes, tally up the voting power of each validator
	if err := k.Votes.Walk(ctx, nil, func(key sdk.AccAddress, vote types.Vote) (bool, error) {
		// if validator, just record it in the map
		voter, err := k.accountKeeper.AddressCodec().StringToBytes(vote.Sender)
		if err != nil {
			return false, err
		}

		// <sunrise>
		// Skip shareclass module's votes
		if sdk.AccAddress(voter).Equals(shareclassAddr) {
			votesToRemove = append(votesToRemove, key)
			return false, nil
		}
		// </sunrise>

		valAddrStr, err := k.stakingKeeper.ValidatorAddressCodec().BytesToString(voter)
		if err != nil {
			return false, err
		}

		if val, ok := validators[valAddrStr]; ok {
			val.PoolWeights = vote.PoolWeights
			validators[valAddrStr] = val
		}

		// iterate over all delegations from voter, deduct from any delegated-to validators
		err = k.stakingKeeper.IterateDelegations(ctx, voter, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
			valAddrStr := delegation.GetValidatorAddr()

			if val, ok := validators[valAddrStr]; ok {
				// There is no need to handle the special case that validator address equal to voter address.
				// Because voter's voting power will tally again even if there will be deduction of voter's voting power from validator.
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				validators[valAddrStr] = val

				// delegation shares * bonded / total shares
				votingPower := delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares)

				for _, poolWeight := range vote.PoolWeights {
					weight, _ := math.LegacyNewDecFromStr(poolWeight.Weight)
					subPower := votingPower.Mul(weight)
					oldWeight := results[poolWeight.PoolId]
					if oldWeight.IsNil() {
						oldWeight = math.LegacyZeroDec()
					}
					results[poolWeight.PoolId] = oldWeight.Add(subPower)
				}
				totalVotingPower = totalVotingPower.Add(votingPower)
			}

			return false
		})
		if err != nil {
			return false, err
		}

		votesToRemove = append(votesToRemove, key)
		return false, nil
	}); err != nil {
		return []types.TallyResult{}, err
	}

	// remove all votes from store
	for _, key := range votesToRemove {
		if err := k.Votes.Remove(ctx, key); err != nil {
			return []types.TallyResult{}, err
		}
	}

	// iterate over the validators again to tally their voting power
	for _, val := range validators {
		if len(val.PoolWeights) == 0 {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

		for _, poolWeight := range val.PoolWeights {
			weight, err := math.LegacyNewDecFromStr(poolWeight.Weight)
			if err != nil {
				return []types.TallyResult{}, err
			}
			subPower := votingPower.Mul(weight)
			oldWeight := results[poolWeight.PoolId]
			if oldWeight.IsNil() {
				oldWeight = math.LegacyZeroDec()
			}
			results[poolWeight.PoolId] = oldWeight.Add(subPower)
		}
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	// If there is no staked coins, the proposal fails
	totalBonded, err := k.stakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		return []types.TallyResult{}, err
	}

	if totalBonded.IsZero() {
		return []types.TallyResult{}, nil
	}

	tallyResults := NewTallyResultFromMap(results)
	sort.SliceStable(tallyResults, func(i, j int) bool {
		return tallyResults[i].PoolId < tallyResults[j].PoolId
	})

	return tallyResults, nil
}

// NewTallyResultFromMap creates a new TallyResult instance from a pool_id -> Dec map
func NewTallyResultFromMap(results map[uint64]math.LegacyDec) []types.TallyResult {
	tallyResults := []types.TallyResult{}
	for poolId, count := range results {
		tallyResults = append(tallyResults, types.TallyResult{
			PoolId: poolId,
			Count:  count.TruncateInt(),
		})
	}
	return tallyResults
}
