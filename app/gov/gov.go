package gov

import (
	"context"

	"cosmossdk.io/collections"
	addresscodec "cosmossdk.io/core/address"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govkeeper "cosmossdk.io/x/gov/keeper"
	v1 "cosmossdk.io/x/gov/types/v1"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

type AccountKeeper interface {
	AddressCodec() addresscodec.Codec

	GetModuleAddress(moduleName string) sdk.AccAddress
}

type StakingKeeper interface {
	ValidatorAddressCodec() addresscodec.Codec

	IterateDelegations(ctx context.Context, delAddr sdk.AccAddress,
		fn func(index int64, del sdk.DelegationI) (stop bool),
	) error

	GetDelegatorBonded(ctx context.Context, delegator sdk.AccAddress) (math.Int, error)
	TotalBondedTokens(ctx context.Context) (math.Int, error)
}

func ProvideCalculateVoteResultsAndVotingPowerFn(authKeeper AccountKeeper, stakingKeeper StakingKeeper) govkeeper.CalculateVoteResultsAndVotingPowerFn {
	return func(
		ctx context.Context,
		keeper govkeeper.Keeper,
		proposalID uint64,
		validators map[string]v1.ValidatorGovInfo,
	) (totalVoterPower math.LegacyDec, results map[v1.VoteOption]math.LegacyDec, err error) {
		totalVP := math.LegacyZeroDec()
		results = createEmptyResults()

		// <sunrise>
		// Deduct shareclass module's delegations
		shareclassAddr := authKeeper.GetModuleAddress(shareclasstypes.ModuleName)
		shareclassVP := math.LegacyZeroDec()
		err = stakingKeeper.IterateDelegations(ctx, shareclassAddr, func(index int64, delegation sdk.DelegationI) (stop bool) {
			valAddrStr := delegation.GetValidatorAddr()
			if val, ok := validators[valAddrStr]; ok {
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				validators[valAddrStr] = val

				shareclassVP = shareclassVP.Add(delegation.GetShares())
			}
			return false
		})
		if err != nil {
			return math.LegacyDec{}, nil, err
		}
		// </sunrise>

		// iterate over all votes, tally up the voting power of each validator
		rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](proposalID)
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
			err = stakingKeeper.IterateDelegations(ctx, voter, func(index int64, delegation sdk.DelegationI) (stop bool) {
				valAddrStr := delegation.GetValidatorAddr()

				if val, ok := validators[valAddrStr]; ok {
					// There is no need to handle the special case that validator address equal to voter address.
					// Because voter's voting power will tally again even if there will be deduction of voter's voting power from validator.
					val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
					validators[valAddrStr] = val

					// delegation shares * bonded / total shares
					votingPower := delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares)

					for _, option := range vote.Options {
						weight, _ := math.LegacyNewDecFromStr(option.Weight)
						subPower := votingPower.Mul(weight)
						results[option.Option] = results[option.Option].Add(subPower)
					}

					totalVP = totalVP.Add(votingPower)
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
			votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

			for _, option := range val.Vote {
				weight, _ := math.LegacyNewDecFromStr(option.Weight)
				subPower := votingPower.Mul(weight)
				results[option.Option] = results[option.Option].Add(subPower)
			}
			totalVP = totalVP.Add(votingPower)
		}

		// <sunrise>
		// To cancel the effect to quorum, we need to adjust the total voting power.
		// totalVoterPowerCustom / totalBonded = (totalVoterPower - shareclassVotingPower) / (totalBonded - shareclassBonded)
		shareclassBonded, err := stakingKeeper.GetDelegatorBonded(ctx, shareclassAddr)
		if err != nil {
			return math.LegacyDec{}, nil, err
		}
		totalBonded, err := stakingKeeper.TotalBondedTokens(ctx)
		if err != nil {
			return math.LegacyDec{}, nil, err
		}
		if !totalBonded.IsZero() {
			numerator := totalVP.Sub(shareclassVP)
			denominator := totalBonded.Sub(shareclassBonded)

			numerator = numerator.MulInt(totalBonded)
			totalVP = numerator.Quo(math.LegacyNewDecFromInt(denominator))
		}
		// <sunrise />

		return totalVP, results, nil
	}
}

func createEmptyResults() map[v1.VoteOption]math.LegacyDec {
	results := make(map[v1.VoteOption]math.LegacyDec)
	results[v1.OptionYes] = math.LegacyZeroDec()
	results[v1.OptionAbstain] = math.LegacyZeroDec()
	results[v1.OptionNo] = math.LegacyZeroDec()
	results[v1.OptionNoWithVeto] = math.LegacyZeroDec()
	results[v1.OptionSpam] = math.LegacyZeroDec()

	return results
}
