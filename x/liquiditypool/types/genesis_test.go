package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
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

				PoolList: []types.Pool{
					{
						Id: 0,
					},
					{
						Id: 1,
					},
				},
				PoolCount: 2,
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated pool",
			genState: &types.GenesisState{
				PoolList: []types.Pool{
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
			desc: "invalid pool count",
			genState: &types.GenesisState{
				PoolList: []types.Pool{
					{
						Id: 1,
					},
				},
				PoolCount: 0,
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
