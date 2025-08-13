package gov

import (
	"context"

	"cosmossdk.io/collections"
	addresscodec "cosmossdk.io/core/address"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

type AccountKeeper interface {
	AddressCodec() addresscodec.Codec

	GetModuleAddress(moduleName string) sdk.AccAddress
}

type StakingKeeper interface {
	ValidatorAddressCodec() addresscodec.Codec

	IterateDelegations(ctx context.Context, delAddr sdk.AccAddress,
		fn func(index int64, del stakingtypes.DelegationI) (stop bool),
	) error

	GetDelegatorBonded(ctx context.Context, delegator sdk.AccAddress) (math.Int, error)
	TotalBondedTokens(ctx context.Context) (math.Int, error)
}

func ProvideCalculateVoteResultsAndVotingPowerFn(authKeeper AccountKeeper, stakingKeeper StakingKeeper) govkeeper.CalculateVoteResultsAndVotingPowerFn {
	return func(
		ctx context.Context,
		keeper govkeeper.Keeper,
		proposal v1.Proposal,
		validators map[string]v1.ValidatorGovInfo,
	) (totalVoterPower math.LegacyDec, results map[v1.VoteOption]math.LegacyDec, err error) {
		totalVotingPower := math.LegacyZeroDec()

		results = make(map[v1.VoteOption]math.LegacyDec)
		results[v1.OptionYes] = math.LegacyZeroDec()
		results[v1.OptionAbstain] = math.LegacyZeroDec()
		results[v1.OptionNo] = math.LegacyZeroDec()
		results[v1.OptionNoWithVeto] = math.LegacyZeroDec()

		// <sunrise>
		// Deduct shareclass module's delegations
		shareclassAddr := authKeeper.GetModuleAddress(shareclasstypes.ModuleName)
		err = stakingKeeper.IterateDelegations(ctx, shareclassAddr, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
			valAddrStr := delegation.GetValidatorAddr()
			if val, ok := validators[valAddrStr]; ok {
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				validators[valAddrStr] = val
			}
			return false
		})
		if err != nil {
			return math.LegacyDec{}, nil, err
		}
		// </sunrise>

		// iterate over all votes, tally up the voting power of each validator
		rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](proposal.Id)
		votesToRemove := []collections.Pair[uint64, sdk.AccAddress]{}
		if err := keeper.Votes.Walk(ctx, rng, func(key collections.Pair[uint64, sdk.AccAddress], vote v1.Vote) (bool, error) {
			// if validator, just record it in the map
			voter, err := authKeeper.AddressCodec().StringToBytes(vote.Voter)
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

			valAddrStr, err := stakingKeeper.ValidatorAddressCodec().BytesToString(voter)
			if err != nil {
				return false, err
			}

			if val, ok := validators[valAddrStr]; ok {
				val.Vote = vote.Options
				validators[valAddrStr] = val
			}

			// iterate over all delegations from voter, deduct from any delegated-to validators
			err = stakingKeeper.IterateDelegations(ctx, voter, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
				valAddrStr := delegation.GetValidatorAddr()

				if val, ok := validators[valAddrStr]; ok {
					// There is no need to handle the special case that validator address equal to voter address.
					// Because voter's voting power will tally again even if there will be deduction of voter's voting power from validator.
					val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
					validators[valAddrStr] = val

					// delegation shares * bonded / total shares
					var votingPower math.LegacyDec
					if val.DelegatorShares.IsZero() {
						votingPower = math.LegacyZeroDec()
					} else {
						votingPower = delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares)
					}

					for _, option := range vote.Options {
						weight, _ := math.LegacyNewDecFromStr(option.Weight)
						subPower := votingPower.Mul(weight)
						results[option.Option] = results[option.Option].Add(subPower)
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
			return math.LegacyDec{}, nil, err
		}

		// remove all votes from store
		for _, key := range votesToRemove {
			if err := keeper.Votes.Remove(ctx, key); err != nil {
				return math.LegacyDec{}, nil, err
			}
		}

		// iterate over the validators again to tally their voting power
		for _, val := range validators {
			if len(val.Vote) == 0 {
				continue
			}

			sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
			var votingPower math.LegacyDec

			if val.DelegatorShares.IsZero() {
				votingPower = math.LegacyZeroDec()
			} else if sharesAfterDeductions.LTE(math.LegacyZeroDec()) {
				votingPower = math.LegacyZeroDec()
			} else {
				votingPower = sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)
			}

			for _, option := range val.Vote {
				weight, _ := math.LegacyNewDecFromStr(option.Weight)
				subPower := votingPower.Mul(weight)
				results[option.Option] = results[option.Option].Add(subPower)
			}
			totalVotingPower = totalVotingPower.Add(votingPower)
		}

		return totalVotingPower, results, nil
	}
}
