package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sunrise-zone/sunrise-app/x/liquidstaking/types"
)

// func TestGenesisState_Validate(t *testing.T) {
// 	tests := []struct {
// 		desc     string
// 		genState *types.GenesisState
// 		valid    bool
// 	}{
// 		{
// 			desc:     "default is valid",
// 			genState: types.DefaultGenesis(),
// 			valid:    true,
// 		},
// 		{
// 			desc:     "valid genesis state",
// 			genState: &types.GenesisState{

// 				// this line is used by starport scaffolding # types/genesis/validField
// 			},
// 			valid: true,
// 		},
// 		// this line is used by starport scaffolding # types/genesis/testcase
// 	}
// 	for _, tc := range tests {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			err := tc.genState.Validate()
// 			if tc.valid {
// 				require.NoError(t, err)
// 			} else {
// 				require.Error(t, err)
// 			}
// 		})
// 	}
// }

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(genState *types.GenesisState)
		expectedErr string
	}{
		{
			"default is valid",
			func(genState *types.GenesisState) {},
			"",
		},
		{
			"invalid liquid validator address",
			func(genState *types.GenesisState) {
				genState.LiquidValidators = []types.LiquidValidator{
					{
						OperatorAddress: "invalidAddr",
					},
				}
			},
			"invalid liquid validator {invalidAddr}: decoding bech32 failed: string not all lowercase or all uppercase: invalid address",
		},
		{
			"empty liquid validator address",
			func(genState *types.GenesisState) {
				genState.LiquidValidators = []types.LiquidValidator{
					{
						OperatorAddress: "",
					},
				}
			},
			"invalid liquid validator {}: empty address string is not allowed: invalid address",
		},
		{
			"invalid params(UnstakeFeeRate)",
			func(genState *types.GenesisState) {
				genState.Params.UnstakeFeeRate = sdk.Dec{}
			},
			"unstake fee rate must not be nil",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			genState := types.DefaultGenesisState()
			tc.malleate(genState)
			err := types.ValidateGenesis(*genState)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
