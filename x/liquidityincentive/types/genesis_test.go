package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Epochs: []types.Epoch{
					{
						Id: 0,
					},
					{
						Id: 1,
					},
				},
				EpochCount: 2,
				Gauges: []types.Gauge{
					{
						PreviousEpochId: 0,
						PoolId:          0,
					},
					{
						PreviousEpochId: 1,
						PoolId:          1,
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated epoch",
			genState: &types.GenesisState{
				Epochs: []types.Epoch{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid epoch count",
			genState: &types.GenesisState{
				Epochs: []types.Epoch{
					{
						Id: 1,
					},
				},
				EpochCount: 0,
			},
			valid: false,
		},
		{
			desc: "duplicated gauge",
			genState: &types.GenesisState{
				Gauges: []types.Gauge{
					{
						PreviousEpochId: 0,
						PoolId:          0,
					},
					{
						PreviousEpochId: 0,
						PoolId:          0,
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
