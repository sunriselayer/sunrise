package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/selfdelegation/keeper"
	"github.com/sunriselayer/sunrise/x/selfdelegation/types"
)

func TestMsgUpdateParams(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	params := types.DefaultParams()
	require.NoError(t, f.keeper.Params.Set(f.ctx, params))

	authorityStr, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	// default params
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
				Params:    params,
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "invalid param",
			input: &types.MsgUpdateParams{
				Authority: authorityStr,
				Params:    types.Params{},
			},
			expErr: true,
		},
		{
			name: "default is valid",
			input: &types.MsgUpdateParams{
				Authority: authorityStr,
				Params:    params,
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
			}
		})
	}
}