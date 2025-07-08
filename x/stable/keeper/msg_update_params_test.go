package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/stable/keeper"
	"github.com/sunriselayer/sunrise/x/stable/types"
)

func TestMsgUpdateParams(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	authorityStr, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	validParams := types.Params{
		AuthorityAddresses: []string{authorityStr},
		StableDenom:        "uusdrise",
	}
	require.NoError(t, validParams.Validate())

	// Set initial params
	require.NoError(t, f.keeper.Params.Set(f.ctx, types.DefaultParams()))

	testCases := []struct {
		name      string
		input     *types.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid authority",
			input: &types.MsgUpdateParams{
				Authority: "invalid",
				Params:    validParams,
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "invalid params",
			input: &types.MsgUpdateParams{
				Authority: authorityStr,
				Params:    types.Params{}, // Invalid because StableDenom is empty
			},
			expErr:    true,
			expErrMsg: "stable denom cannot be empty",
		},
		{
			name: "all good",
			input: &types.MsgUpdateParams{
				Authority: authorityStr,
				Params:    validParams,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.UpdateParams(f.ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
				// Check if params were updated
				params, err := f.keeper.Params.Get(f.ctx)
				require.NoError(t, err)
				require.Equal(t, validParams, params)
			}
		})
	}
}
