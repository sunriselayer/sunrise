package keeper

import (
	"context"
	"sort"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
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
	results := make(map[uint64]math.LegacyDec)

	totalVotingPower := math.LegacyZeroDec()
	currValidators := make(map[string]ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	err := k.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		valBz, err := k.sk.ValidatorAddressCodec().StringToBytes(validator.GetOperator())
		if err != nil {
			return false
		}
		currValidators[validator.GetOperator()] = ValidatorGovInfo{
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

	votes := k.GetAllVotes(ctx)
	for _, vote := range votes {
		// if validator, just record it in the map
		voter, err := k.authKeeper.AddressCodec().StringToBytes(vote.Sender)
		if err != nil {
			return []types.TallyResult{}, err
		}

		valAddrStr, err := k.sk.ValidatorAddressCodec().BytesToString(voter)
		if err != nil {
			return []types.TallyResult{}, err
		}
		if val, ok := currValidators[valAddrStr]; ok {
			val.PoolWeights = vote.PoolWeights
			currValidators[valAddrStr] = val
		}

		// iterate over all delegations from voter, deduct from any delegated-to validators
		err = k.sk.IterateDelegations(ctx, voter, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
			valAddrStr := delegation.GetValidatorAddr()

			if val, ok := currValidators[valAddrStr]; ok {
				// There is no need to handle the special case that validator address equal to voter address.
				// Because voter's voting power will tally again even if there will be deduction of voter's voting power from validator.
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				currValidators[valAddrStr] = val

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
			return []types.TallyResult{}, err
		}
	}

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if len(val.PoolWeights) == 0 {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

		for _, poolWeight := range val.PoolWeights {
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

	// If there is no staked coins, the proposal fails
	totalBonded, err := k.sk.TotalBondedTokens(ctx)
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
