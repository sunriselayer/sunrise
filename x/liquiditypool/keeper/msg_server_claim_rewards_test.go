package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestMsgServerClaimRewards(t *testing.T) {
	initFixture(t)
	sender := sdk.AccAddress("sender")

	tests := []struct {
		desc       string
		request    *types.MsgClaimRewards
		allocation sdk.Coins
		expRewards sdk.Coins
		err        error
	}{
		{
			desc: "Single token rewards",
			request: &types.MsgClaimRewards{
				Sender:      sender.String(),
				PositionIds: []uint64{0},
			},
			allocation: sdk.Coins{sdk.NewInt64Coin("xyz", 1000)},
			expRewards: sdk.Coins{sdk.NewInt64Coin("xyz", 999)},
		},
		{
			desc: "Multiple token rewards",
			request: &types.MsgClaimRewards{
				Sender:      sender.String(),
				PositionIds: []uint64{0},
			},
			allocation: sdk.Coins{sdk.NewInt64Coin("uvw", 1000), sdk.NewInt64Coin("xyz", 1000)},
			expRewards: sdk.Coins{sdk.NewInt64Coin("uvw", 999), sdk.NewInt64Coin("xyz", 999)},
		},
		{
			desc: "Empty rewards",
			request: &types.MsgClaimRewards{
				Sender:      sender.String(),
				PositionIds: []uint64{1},
			},
			expRewards: sdk.Coins{},
		},
		{
			desc: "Not available position",
			request: &types.MsgClaimRewards{
				Sender:      sender.String(),
				PositionIds: []uint64{3},
			},
			expRewards: sdk.Coins{},
			err:        types.ErrPositionNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, mocks, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			_, err := srv.CreatePool(wctx, &types.MsgCreatePool{
				Authority:  sender.String(),
				DenomBase:  "base",
				DenomQuote: "quote",
				FeeRate:    "0.01",
				PriceRatio: "1.0001",
				BaseOffset: "0.5",
			})
			require.NoError(t, err)

			_, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
				Sender:         sender.String(),
				PoolId:         0,
				LowerTick:      -10,
				UpperTick:      10,
				TokenBase:      sdk.NewInt64Coin("base", 10000),
				TokenQuote:     sdk.NewInt64Coin("quote", 10000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			_, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
				Sender:         sender.String(),
				PoolId:         0,
				LowerTick:      100,
				UpperTick:      200,
				TokenBase:      sdk.NewInt64Coin("base", 10000),
				TokenQuote:     sdk.NewInt64Coin("quote", 10000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			err = k.AllocateIncentive(wctx, 0, sender, tc.allocation)
			require.NoError(t, err)

			resp, err := srv.ClaimRewards(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, resp.ClaimedFees.String(), tc.expRewards.String())
			}
		})
	}
}
