package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/stable/types"
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
			desc: "invalid genesis state - invalid authority address",
			genState: &types.GenesisState{
				Params: types.Params{
					AuthorityAddresses: []string{"invalid-address"},
					StableDenom:        "uusdrise",
				},
			},
			valid: false,
		},
		{
			desc: "invalid genesis state - empty stable denom",
			genState: &types.GenesisState{
				Params: types.Params{
					AuthorityAddresses: []string{"cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqn0kn67"},
					StableDenom:        "",
				},
			},
			valid: false,
		},
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
