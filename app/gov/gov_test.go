package gov_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// This function is a direct copy of the logic from `gov.go`, modified to be a standalone, testable function.
func calculate(
	t *testing.T,
	totalBonded math.Int,
	shareclassBonded math.Int,
	delegations map[string][]stakingtypes.Delegation,
	votes []v1.Vote,
	validators map[string]v1.ValidatorGovInfo,
) (math.LegacyDec, map[v1.VoteOption]math.LegacyDec) {
	shareclassAddr := sdk.AccAddress("shareclass_module")

	totalVotingPower := math.LegacyZeroDec()
	results := make(map[v1.VoteOption]math.LegacyDec)
	results[v1.OptionYes] = math.LegacyZeroDec()
	results[v1.OptionAbstain] = math.LegacyZeroDec()
	results[v1.OptionNo] = math.LegacyZeroDec()
	results[v1.OptionNoWithVeto] = math.LegacyZeroDec()

	// 1. Deduct shareclass module's delegations
	for _, delegation := range delegations[shareclassAddr.String()] {
		valAddrStr := delegation.GetValidatorAddr()
		if val, ok := validators[valAddrStr]; ok {
			val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
			validators[valAddrStr] = val
		}
	}

	// 2. Tally delegator votes
	for _, vote := range votes {
		voter, err := sdk.AccAddressFromBech32(vote.Voter)
		require.NoError(t, err)

		if voter.Equals(shareclassAddr) {
			continue // Skip shareclass module's votes
		}

		voterDels, ok := delegations[voter.String()]
		if !ok {
			continue
		}

		for _, delegation := range voterDels {
			valAddrStr := delegation.GetValidatorAddr()
			if val, ok := validators[valAddrStr]; ok {
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				validators[valAddrStr] = val

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
		}
	}

	// 3. Tally validator votes
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

	// 4. Scale the results
	effectiveTotalBonded := totalBonded.Sub(shareclassBonded)
	if !effectiveTotalBonded.IsPositive() || totalBonded.IsZero() {
		return totalVotingPower, results
	}

	for i, r := range results {
		results[i] = r.MulInt(totalBonded).QuoInt(effectiveTotalBonded)
	}
	totalVotingPower = totalVotingPower.MulInt(totalBonded).QuoInt(effectiveTotalBonded)

	return totalVotingPower, results
}

func TestTallyLogic(t *testing.T) {
	valAddr1 := sdk.ValAddress("validator1")
	valAddrStr1, _ := sdk.Bech32ifyAddressBytes("cosmosvaloper", valAddr1)
	shareclassAddr := sdk.AccAddress("shareclass_module")
	delegatorAddr := sdk.AccAddress("delegator1")

	testCases := []struct {
		name               string
		totalBonded        math.Int
		shareclassBonded   math.Int
		delegations        map[string][]stakingtypes.Delegation
		votes              []v1.Vote
		validators         map[string]v1.ValidatorGovInfo
		expectedTotalPower math.LegacyDec
		expectedResults    map[v1.VoteOption]math.LegacyDec
	}{
		{
			name:             "Scenario: Validator and Delegator vote, with shareclass",
			totalBonded:      math.NewInt(1000),
			shareclassBonded: math.NewInt(100),
			delegations: map[string][]stakingtypes.Delegation{
				shareclassAddr.String(): {
					stakingtypes.NewDelegation(shareclassAddr.String(), valAddr1.String(), math.LegacyNewDec(100)),
				},
				delegatorAddr.String(): {
					stakingtypes.NewDelegation(delegatorAddr.String(), valAddr1.String(), math.LegacyNewDec(200)),
				},
			},
			votes: []v1.Vote{
				{Voter: delegatorAddr.String(), Options: v1.NewNonSplitVoteOption(v1.OptionNo)},
			},
			validators: map[string]v1.ValidatorGovInfo{
				valAddrStr1: {
					Vote:            v1.NewNonSplitVoteOption(v1.OptionYes),
					BondedTokens:    math.NewInt(1000),
					DelegatorShares: math.LegacyNewDec(1000),
				},
			},
			// Delegator (No): power = 200
			// Validator (Yes): power = 1000(total) - 100(shareclass) - 200(delegator) = 700
			// Total Power = 200 + 700 = 900
			// Scale Factor: 1000 / (1000 - 100) = 1000 / 900
			expectedTotalPower: math.LegacyNewDec(900).MulInt(math.NewInt(1000)).QuoInt(math.NewInt(900)), // Should be 900
			expectedResults: map[v1.VoteOption]math.LegacyDec{
				v1.OptionYes:        math.LegacyNewDec(700).MulInt(math.NewInt(1000)).QuoInt(math.NewInt(900)),
				v1.OptionAbstain:    math.LegacyZeroDec(),
				v1.OptionNo:         math.LegacyNewDec(200).MulInt(math.NewInt(1000)).QuoInt(math.NewInt(900)),
				v1.OptionNoWithVeto: math.LegacyZeroDec(),
			},
		},
		{
			name:             "Scenario: Shareclass module vote is ignored",
			totalBonded:      math.NewInt(1000),
			shareclassBonded: math.ZeroInt(),
			delegations:      map[string][]stakingtypes.Delegation{},
			votes: []v1.Vote{
				{Voter: shareclassAddr.String(), Options: v1.NewNonSplitVoteOption(v1.OptionYes)},
			},
			validators: map[string]v1.ValidatorGovInfo{
				valAddrStr1: {BondedTokens: math.NewInt(1000), DelegatorShares: math.LegacyNewDec(1000)},
			},
			expectedTotalPower: math.LegacyZeroDec(), // No one with power voted
			expectedResults: map[v1.VoteOption]math.LegacyDec{
				v1.OptionYes:        math.LegacyZeroDec(),
				v1.OptionAbstain:    math.LegacyZeroDec(),
				v1.OptionNo:         math.LegacyZeroDec(),
				v1.OptionNoWithVeto: math.LegacyZeroDec(),
			},
		},
		{
			name:             "Edge Case: Shareclass bonded equals total bonded",
			totalBonded:      math.NewInt(1000),
			shareclassBonded: math.NewInt(1000),
			delegations: map[string][]stakingtypes.Delegation{
				shareclassAddr.String(): {
					stakingtypes.NewDelegation(shareclassAddr.String(), valAddr1.String(), math.LegacyNewDec(1000)),
				},
			},
			votes: []v1.Vote{},
			validators: map[string]v1.ValidatorGovInfo{
				valAddrStr1: {
					Vote:            v1.NewNonSplitVoteOption(v1.OptionYes),
					BondedTokens:    math.NewInt(1000),
					DelegatorShares: math.LegacyNewDec(1000),
				},
			},
			expectedTotalPower: math.LegacyZeroDec(), // Power becomes 0, no scaling
			expectedResults: map[v1.VoteOption]math.LegacyDec{
				v1.OptionYes:        math.LegacyZeroDec(),
				v1.OptionAbstain:    math.LegacyZeroDec(),
				v1.OptionNo:         math.LegacyZeroDec(),
				v1.OptionNoWithVeto: math.LegacyZeroDec(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validatorsCopy := make(map[string]v1.ValidatorGovInfo)
			for k, v := range tc.validators {
				v.DelegatorDeductions = math.LegacyZeroDec()
				validatorsCopy[k] = v
			}

			totalPower, results := calculate(
				t,
				tc.totalBonded,
				tc.shareclassBonded,
				tc.delegations,
				tc.votes,
				validatorsCopy,
			)

			require.True(t, tc.expectedTotalPower.Equal(totalPower), "expected total power %s, got %s", tc.expectedTotalPower, totalPower)
			for option, expected := range tc.expectedResults {
				require.True(t, expected.Equal(results[option]), "results for option %s: expected %s, got %s", option, expected, results[option])
			}
		})
	}
}
