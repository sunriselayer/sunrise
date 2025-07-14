package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/x/fee/keeper"
	"github.com/sunriselayer/sunrise/x/fee/types"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestMsgUpdateParams(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	params := types.DefaultParams()
	require.NoError(t, f.keeper.Params.Set(f.ctx, params))

	authorityStr := f.keeper.GetAuthority()

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
			name: "burn enabled with invalid pool",
			input: &types.MsgUpdateParams{
				Authority: authorityStr,
				Params:    types.NewParams("fee", "burn", math.LegacyNewDecWithPrec(50, 2), 1, true),
			},
			expErr:    true,
			expErrMsg: "invalid pool",
		},
		{
			name: "burn enabled with valid pool",
			input: &types.MsgUpdateParams{
				Authority: authorityStr,
				Params:    types.NewParams(consts.StableDenom, consts.MintDenom, math.LegacyNewDecWithPrec(50, 2), 1, true),
			},
			expErr: false,
		},
		{
			name: "burn disabled",
			input: &types.MsgUpdateParams{
				Authority: authorityStr,
				Params:    params,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f.mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), gomock.Any()).Return(liquiditypooltypes.Pool{
				Id:         1,
				DenomBase:  consts.StableDenom,
				DenomQuote: consts.MintDenom,
			}, true, nil).AnyTimes()
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
