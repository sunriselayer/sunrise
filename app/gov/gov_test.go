package gov_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// This function is a direct copy of the logic from `gov.go`, modified to be
// a standalone, testable function.
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

	// Deduct shareclass module's delegations
	for _, delegation := range delegations[shareclassAddr.String()] {
		valAddrStr := delegation.GetValidatorAddr()
		if val, ok := validators[valAddrStr]; ok {
			val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
			validators[valAddrStr] = val
		}
	}

	// Tally delegator votes
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

	// Tally validator votes
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

	// Scale the results
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
			name:             "Scenario: Basic shareclass exclusion and delegator vote",
			totalBonded:      math.NewInt(2000),
			shareclassBonded: math.NewInt(500),
			delegations: map[string][]stakingtypes.Delegation{
				shareclassAddr.String(): {
					stakingtypes.NewDelegation(shareclassAddr.String(), valAddr1.String(), math.LegacyNewDec(500)),
				},
				delegatorAddr.String(): {
					stakingtypes.NewDelegation(delegatorAddr.String(), valAddr1.String(), math.LegacyNewDec(100)),
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
			// Delegator(No): power = 100
			// Validator(Yes): power = (1000(total) - 500(shareclass) - 100(delegator)) = 400
			// Total Power = 100 + 400 = 500
			// Scale Factor: 2000 / (2000 - 500) = 1.333...
			// Expected Power = 500 * 1.333 = 666.66...
			// Expected Yes = 400 * 1.333 = 533.33...
			// Expected No = 100 * 1.333 = 133.33...
			// Let's use fractions for precision: 2000/1500 = 4/3
			// Expected Power = 500 * 4/3 = 666.666...
			expectedTotalPower: math.LegacyNewDec(500).MulInt(math.NewInt(2000)).QuoInt(math.NewInt(1500)),
			expectedResults: map[v1.VoteOption]math.LegacyDec{
				v1.OptionYes:        math.LegacyNewDec(400).MulInt(math.NewInt(2000)).QuoInt(math.NewInt(1500)),
				v1.OptionAbstain:    math.LegacyZeroDec(),
				v1.OptionNo:         math.LegacyNewDec(100).MulInt(math.NewInt(2000)).QuoInt(math.NewInt(1500)),
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
			expectedTotalPower: math.LegacyZeroDec(), // No one voted
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
			validators: map[string]v1.ValidatorGovInfo{
				valAddrStr1: {
					Vote:            v1.NewNonSplitVoteOption(v1.OptionYes),
					BondedTokens:    math.NewInt(1000),
					DelegatorShares: math.LegacyNewDec(1000),
				},
			},
			expectedTotalPower: math.LegacyZeroDec(), // Power is 0, no scaling
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
			// Make a deep copy of validators map to avoid modifying the original
			validatorsCopy := make(map[string]v1.ValidatorGovInfo)
			for k, v := range tc.validators {
				// Initialize deductions for the copy
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
