package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	"github.com/sunrise-zone/sunrise-app/x/liquidstaking/types"
)

var (
	liquidValidators = []types.LiquidValidator{
		{
			OperatorAddress: "cosmosvaloper15kdfwczhpmccprekhlzrvkhzw92940l3w37qqj",
		},
		{
			OperatorAddress: "cosmosvaloper1x73gyvh74ahs2rt9cqrpjkkk74nczwfpnskv3rczmsf0m6aj5dksqr58m3",
		},
		{
			OperatorAddress: "cosmosvaloper10ngyx42lfpylpllm4k3g7fz4gufnt3ptyhm5pn",
		},
		{
			OperatorAddress: "cosmosvaloper10fcwju2n8vvffkp8judj3skqpvnphasxjar5yx",
		},
	}
)

func TestDivideByWeight(t *testing.T) {
	testCases := []struct {
		whitelistedVals  []types.WhitelistedValidator
		addStakingAmt    sdkmath.Int
		currentDelShares []sdkmath.Int
		expectedOutputs  []sdkmath.Int
		expectedCrumb    sdkmath.Int
	}{
		{
			whitelistedVals: []types.WhitelistedValidator{
				{
					ValidatorAddress: liquidValidators[0].OperatorAddress,
					TargetWeight:     sdkmath.NewInt(1),
				},
				{
					ValidatorAddress: liquidValidators[1].OperatorAddress,
					TargetWeight:     sdkmath.NewInt(1),
				},
				{
					ValidatorAddress: liquidValidators[2].OperatorAddress,
					TargetWeight:     sdkmath.NewInt(1),
				},
			},
			addStakingAmt:    sdkmath.NewInt(10 * 1000000),
			currentDelShares: []sdkmath.Int{sdkmath.NewInt(2000000), sdkmath.NewInt(2000000), sdkmath.NewInt(1000000)},
			expectedOutputs:  []sdkmath.Int{sdkmath.NewInt(3333333), sdkmath.NewInt(3333333), sdkmath.NewInt(3333333)},
			expectedCrumb:    sdkmath.NewInt(1),
		},
		{
			whitelistedVals: []types.WhitelistedValidator{
				{
					ValidatorAddress: liquidValidators[0].OperatorAddress,
					TargetWeight:     sdkmath.NewInt(2),
				},
				{
					ValidatorAddress: liquidValidators[1].OperatorAddress,
					TargetWeight:     sdkmath.NewInt(2),
				},
				{
					ValidatorAddress: liquidValidators[2].OperatorAddress,
					TargetWeight:     sdkmath.NewInt(1),
				},
			},
			addStakingAmt:    sdkmath.NewInt(10 * 1000000),
			currentDelShares: []sdkmath.Int{sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.NewInt(1000000)},
			expectedOutputs:  []sdkmath.Int{sdkmath.NewInt(4000000), sdkmath.NewInt(4000000), sdkmath.NewInt(2000000)},
			expectedCrumb:    sdkmath.NewInt(0),
		},
		{
			whitelistedVals: []types.WhitelistedValidator{
				{
					ValidatorAddress: liquidValidators[0].OperatorAddress,
					TargetWeight:     sdkmath.NewInt(1),
				},
				{
					ValidatorAddress: liquidValidators[1].OperatorAddress,
					TargetWeight:     sdkmath.NewInt(1),
				},
				{
					ValidatorAddress: liquidValidators[2].OperatorAddress,
					TargetWeight:     sdkmath.NewInt(1),
				},
			},
			addStakingAmt:    sdkmath.NewInt(10),
			currentDelShares: []sdkmath.Int{sdkmath.NewInt(3), sdkmath.NewInt(2), sdkmath.NewInt(1)},
			expectedOutputs:  []sdkmath.Int{sdkmath.NewInt(3), sdkmath.NewInt(3), sdkmath.NewInt(3)},
			expectedCrumb:    sdkmath.NewInt(1),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, []types.WhitelistedValidator{}, tc.whitelistedVals)
		require.IsType(t, sdkmath.Int{}, tc.addStakingAmt)
		require.IsType(t, sdkmath.Int{}, tc.expectedCrumb)
		require.IsType(t, []sdkmath.Int{}, tc.expectedOutputs)

		totalTargetAmt := sdkmath.ZeroInt()
		valsMap := types.GetWhitelistedValsMap(tc.whitelistedVals)
		var activeVals types.ActiveLiquidValidators
		for _, v := range tc.whitelistedVals {
			activeVals = append(activeVals, types.LiquidValidator{
				OperatorAddress: v.ValidatorAddress,
			})
		}
		outputs, crumb := types.DivideByWeight(activeVals, tc.addStakingAmt, valsMap)
		for _, v := range outputs {
			totalTargetAmt = totalTargetAmt.Add(v)
		}
		require.EqualValues(t, tc.expectedOutputs, outputs)
		require.EqualValues(t, tc.addStakingAmt, totalTargetAmt.Add(crumb))
		require.Equal(t, tc.expectedCrumb.String(), crumb.String())
	}
}

func TestMinMaxGap(t *testing.T) {
	testCases := []struct {
		name                     string
		liquidVals               types.LiquidValidators
		targetMap                map[string]sdkmath.Int
		liquidTokenMap           map[string]sdkmath.Int
		expectedMinGapVal        types.LiquidValidator
		expectedMaxGapVal        types.LiquidValidator
		expectedAmountNeeded     sdkmath.Int
		expectedLastRedelegation bool
	}{
		{
			name:       "zero case",
			liquidVals: liquidValidators,
			targetMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.ZeroInt(),
				liquidValidators[1].OperatorAddress: sdkmath.ZeroInt(),
				liquidValidators[2].OperatorAddress: sdkmath.ZeroInt(),
				liquidValidators[3].OperatorAddress: sdkmath.ZeroInt(),
			},
			liquidTokenMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.ZeroInt(),
				liquidValidators[1].OperatorAddress: sdkmath.ZeroInt(),
				liquidValidators[2].OperatorAddress: sdkmath.ZeroInt(),
				liquidValidators[3].OperatorAddress: sdkmath.ZeroInt(),
			},
			expectedMinGapVal:        types.LiquidValidator{},
			expectedMaxGapVal:        types.LiquidValidator{},
			expectedAmountNeeded:     sdkmath.ZeroInt(),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 1-1",
			liquidVals: liquidValidators,
			targetMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[3].OperatorAddress: sdkmath.NewInt(100000000),
			},
			liquidTokenMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(133333334),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[3].OperatorAddress: sdkmath.ZeroInt(),
			},
			expectedMinGapVal:        liquidValidators[3],
			expectedMaxGapVal:        liquidValidators[0],
			expectedAmountNeeded:     sdkmath.NewInt(33333334),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 1-2",
			liquidVals: liquidValidators,
			targetMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[3].OperatorAddress: sdkmath.NewInt(100000000),
			},
			liquidTokenMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(133333334 - 33333334),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[3].OperatorAddress: sdkmath.NewInt(0 + 33333334),
			},
			expectedMinGapVal:        liquidValidators[3],
			expectedMaxGapVal:        liquidValidators[1],
			expectedAmountNeeded:     sdkmath.NewInt(33333333),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 1-3",
			liquidVals: liquidValidators,
			targetMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[3].OperatorAddress: sdkmath.NewInt(100000000),
			},
			liquidTokenMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(133333334 - 33333334),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(133333333 - 33333333),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[3].OperatorAddress: sdkmath.NewInt(33333334 + 33333333),
			},
			expectedMinGapVal:        liquidValidators[3],
			expectedMaxGapVal:        liquidValidators[2],
			expectedAmountNeeded:     sdkmath.NewInt(33333333),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 1-4",
			liquidVals: liquidValidators,
			targetMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[3].OperatorAddress: sdkmath.NewInt(100000000),
			},
			liquidTokenMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(133333334 - 33333334),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(133333333 - 33333333),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(133333333 - 33333333),
				liquidValidators[3].OperatorAddress: sdkmath.NewInt(33333334 + 33333333 + 33333333),
			},
			expectedMinGapVal:        types.LiquidValidator{},
			expectedMaxGapVal:        types.LiquidValidator{},
			expectedAmountNeeded:     sdkmath.ZeroInt(),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 2-1",
			liquidVals: liquidValidators,
			targetMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(133333334),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[3].OperatorAddress: sdkmath.ZeroInt(),
			},
			liquidTokenMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[3].OperatorAddress: sdkmath.NewInt(100000000),
			},
			expectedMinGapVal:        liquidValidators[0],
			expectedMaxGapVal:        liquidValidators[3],
			expectedAmountNeeded:     sdkmath.NewInt(33333334),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 2-2",
			liquidVals: liquidValidators,
			targetMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(133333334),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[3].OperatorAddress: sdkmath.ZeroInt(),
			},
			liquidTokenMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(100000000 + 33333334),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[3].OperatorAddress: sdkmath.NewInt(100000000 - 33333334),
			},
			expectedMinGapVal:        liquidValidators[1],
			expectedMaxGapVal:        liquidValidators[3],
			expectedAmountNeeded:     sdkmath.NewInt(33333333),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 2-3, last redelegation",
			liquidVals: liquidValidators,
			targetMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(133333334),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(133333333),
				liquidValidators[3].OperatorAddress: sdkmath.ZeroInt(),
			},
			liquidTokenMap: map[string]sdkmath.Int{
				liquidValidators[0].OperatorAddress: sdkmath.NewInt(100000000 + 33333334),
				liquidValidators[1].OperatorAddress: sdkmath.NewInt(100000000 + 33333333),
				liquidValidators[2].OperatorAddress: sdkmath.NewInt(100000000),
				liquidValidators[3].OperatorAddress: sdkmath.NewInt(100000000 - 33333334 - 33333333),
			},
			expectedMinGapVal:        liquidValidators[2],
			expectedMaxGapVal:        liquidValidators[3],
			expectedAmountNeeded:     sdkmath.NewInt(33333333),
			expectedLastRedelegation: true,
		},
	}

	for _, tc := range testCases {
		minGapVal, maxGapVal, amountNeeded, last := tc.liquidVals.MinMaxGap(tc.targetMap, tc.liquidTokenMap)
		require.EqualValues(t, minGapVal, tc.expectedMinGapVal)
		require.EqualValues(t, maxGapVal, tc.expectedMaxGapVal)
		require.EqualValues(t, amountNeeded, tc.expectedAmountNeeded)
		require.EqualValues(t, last, tc.expectedLastRedelegation)
	}
}

func TestDivideByCurrentWeight(t *testing.T) {
	testCases := []struct {
		liquidValidators []types.LiquidValidatorState
		addStakingAmt    sdkmath.LegacyDec
		expectedOutputs  []sdkmath.LegacyDec
		expectedCrumb    sdkmath.LegacyDec
	}{
		{
			liquidValidators: []types.LiquidValidatorState{
				{
					OperatorAddress: "a",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(2 * 1000000),
				},
				{
					OperatorAddress: "b",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(2 * 1000000),
				},
				{
					OperatorAddress: "c",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(1 * 1000000),
				},
			},
			addStakingAmt:   sdkmath.LegacyNewDec(10 * 1000000),
			expectedOutputs: []sdkmath.LegacyDec{sdkmath.LegacyNewDec(4 * 1000000), sdkmath.LegacyNewDec(4 * 1000000), sdkmath.LegacyNewDec(2 * 1000000)},
			expectedCrumb:   sdkmath.LegacyNewDec(0),
		},
		{
			liquidValidators: []types.LiquidValidatorState{
				{
					OperatorAddress: "a",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(1 * 1000000),
					Weight:          sdkmath.NewInt(2),
				},
				{
					OperatorAddress: "b",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(1 * 1000000),
					Weight:          sdkmath.NewInt(2),
				},
				{
					OperatorAddress: "c",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(1 * 1000000),
					Weight:          sdkmath.NewInt(1),
				},
			},
			addStakingAmt:   sdkmath.LegacyNewDec(10 * 1000000),
			expectedOutputs: []sdkmath.LegacyDec{sdkmath.LegacyMustNewDecFromStr("3333333.000000000000000000"), sdkmath.LegacyMustNewDecFromStr("3333333.000000000000000000"), sdkmath.LegacyMustNewDecFromStr("3333333.000000000000000000")},
			expectedCrumb:   sdkmath.LegacyMustNewDecFromStr("1.000000000000000000"),
		},
		{
			liquidValidators: []types.LiquidValidatorState{
				{
					OperatorAddress: "a",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(3),
				},
				{
					OperatorAddress: "b",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(2),
				},
				{
					OperatorAddress: "c",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(1),
				},
			},
			addStakingAmt:   sdkmath.LegacyNewDec(10),
			expectedOutputs: []sdkmath.LegacyDec{sdkmath.LegacyMustNewDecFromStr("4.000000000000000000"), sdkmath.LegacyMustNewDecFromStr("3.000000000000000000"), sdkmath.LegacyMustNewDecFromStr("1.000000000000000000")},
			expectedCrumb:   sdkmath.LegacyMustNewDecFromStr("2.000000000000000000"),
		},
		{
			liquidValidators: []types.LiquidValidatorState{
				{
					OperatorAddress: "a",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(10000000),
				},
				{
					OperatorAddress: "b",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(2000000),
				},
				{
					OperatorAddress: "c",
					Status:          types.ValidatorStatus_ValidatorStatusActive,
					LiquidTokens:    sdkmath.NewIntFromUint64(3000001),
				},
			},
			addStakingAmt:   sdkmath.LegacyNewDec(10000000),
			expectedOutputs: []sdkmath.LegacyDec{sdkmath.LegacyMustNewDecFromStr("6666666.000000000000000000"), sdkmath.LegacyMustNewDecFromStr("1333333.000000000000000000"), sdkmath.LegacyMustNewDecFromStr("2000000.000000000000000000")},
			expectedCrumb:   sdkmath.LegacyMustNewDecFromStr("1.000000000000000000"),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, []types.LiquidValidatorState{}, tc.liquidValidators)
		require.IsType(t, sdkmath.LegacyDec{}, tc.addStakingAmt)
		require.IsType(t, sdkmath.LegacyDec{}, tc.expectedCrumb)
		require.IsType(t, []sdkmath.LegacyDec{}, tc.expectedOutputs)

		totalTargetAmt := sdkmath.LegacyZeroDec()
		totalLiquidTokens := sdkmath.ZeroInt()
		liquidTokenMap := map[string]sdkmath.Int{}
		var lvs types.LiquidValidators
		for _, v := range tc.liquidValidators {
			totalLiquidTokens = totalLiquidTokens.Add(v.LiquidTokens)
			liquidTokenMap[v.OperatorAddress] = v.LiquidTokens
			lvs = append(lvs, types.LiquidValidator{
				OperatorAddress: v.OperatorAddress})
		}
		outputs, crumb := types.DivideByCurrentWeight(lvs, tc.addStakingAmt, totalLiquidTokens, liquidTokenMap)
		for _, v := range outputs {
			totalTargetAmt = totalTargetAmt.Add(v)
		}
		require.EqualValues(t, tc.expectedOutputs, outputs)
		require.EqualValues(t, tc.addStakingAmt, totalTargetAmt.Add(crumb))
		require.Equal(t, tc.expectedCrumb.String(), crumb.String())
	}
}
