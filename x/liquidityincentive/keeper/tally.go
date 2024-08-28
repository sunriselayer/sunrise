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
	Weights             []types.PoolWeight // Vote of the validator
}

func (k Keeper) Tally(ctx context.Context) ([]types.PoolWeight, error) {
	weights := make(map[uint64]math.LegacyDec)

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
			Weights:             []types.PoolWeight{},
		}

		return false
	})
	if err != nil {
		return []types.PoolWeight{}, err
	}

	votes := k.GetAllVotes(ctx)
	for _, vote := range votes {
		// if validator, just record it in the map
		voter, err := k.authKeeper.AddressCodec().StringToBytes(vote.Sender)
		if err != nil {
			return []types.PoolWeight{}, err
		}

		valAddrStr, err := k.sk.ValidatorAddressCodec().BytesToString(voter)
		if err != nil {
			return []types.PoolWeight{}, err
		}
		if val, ok := currValidators[valAddrStr]; ok {
			val.Weights = vote.Weights
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

				for _, weight := range vote.Weights {
					subPower := votingPower.Mul(weight.Weight)
					oldWeight := weights[weight.PoolId]
					if oldWeight.IsNil() {
						oldWeight = math.LegacyZeroDec()
					}
					weights[weight.PoolId] = oldWeight.Add(subPower)
				}
				totalVotingPower = totalVotingPower.Add(votingPower)
			}

			return false
		})
		if err != nil {
			return []types.PoolWeight{}, err
		}
	}

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if len(val.Weights) == 0 {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

		for _, weight := range val.Weights {
			subPower := votingPower.Mul(weight.Weight)
			oldWeight := weights[weight.PoolId]
			if oldWeight.IsNil() {
				oldWeight = math.LegacyZeroDec()
			}
			weights[weight.PoolId] = oldWeight.Add(subPower)
		}
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	// If there is no staked coins, the proposal fails
	totalBonded, err := k.sk.TotalBondedTokens(ctx)
	if err != nil {
		return []types.PoolWeight{}, err
	}

	if totalBonded.IsZero() {
		return []types.PoolWeight{}, nil
	}

	weightsArr := []types.PoolWeight{}
	for poolId, weight := range weights {
		weightsArr = append(weightsArr, types.PoolWeight{
			PoolId: poolId,
			Weight: weight,
		})
	}
	sort.SliceStable(weightsArr, func(i, j int) bool {
		return weightsArr[i].PoolId < weightsArr[j].PoolId
	})

	return weightsArr, nil
}
