package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestMsgServerClaimRewards(t *testing.T) {
	sender := "sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f"
	senderAcc := sdk.MustAccAddressFromBech32(sender)

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
				Sender:      sender,
				PositionIds: []uint64{1},
			},
			allocation: sdk.Coins{sdk.NewInt64Coin("xyz", 1000)},
			expRewards: sdk.Coins{sdk.NewInt64Coin("xyz", 999)},
		},
		{
			desc: "Multiple token rewards",
			request: &types.MsgClaimRewards{
				Sender:      sender,
				PositionIds: []uint64{1},
			},
			allocation: sdk.Coins{sdk.NewInt64Coin("uvw", 1000), sdk.NewInt64Coin("xyz", 1000)},
			expRewards: sdk.Coins{sdk.NewInt64Coin("uvw", 999), sdk.NewInt64Coin("xyz", 999)},
		},
		{
			desc: "Empty rewards",
			request: &types.MsgClaimRewards{
				Sender:      sender,
				PositionIds: []uint64{2},
			},
			expRewards: sdk.Coins{},
		},
		{
			desc: "Not available position",
			request: &types.MsgClaimRewards{
				Sender:      sender,
				PositionIds: []uint64{4},
			},
			expRewards: sdk.Coins{},
			err:        types.ErrPositionNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, bk, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			bk.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			bk.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			_, err := srv.CreatePool(wctx, &types.MsgCreatePool{
				Authority:  sender,
				DenomBase:  "base",
				DenomQuote: "quote",
				FeeRate:    "0.01",
				PriceRatio: "1.0001",
				BaseOffset: "0.5",
			})
			require.NoError(t, err)

			_, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
				Sender:         sender,
				PoolId:         1,
				LowerTick:      -10,
				UpperTick:      10,
				TokenBase:      sdk.NewInt64Coin("base", 10000),
				TokenQuote:     sdk.NewInt64Coin("quote", 10000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			_, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
				Sender:         sender,
				PoolId:         1,
				LowerTick:      100,
				UpperTick:      200,
				TokenBase:      sdk.NewInt64Coin("base", 10000),
				TokenQuote:     sdk.NewInt64Coin("quote", 10000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			err = k.AllocateIncentive(wctx, 1, senderAcc, tc.allocation)
			require.NoError(t, err)

			resp, err := srv.ClaimRewards(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, resp.CollectedFees.String(), tc.expRewards.String())
			}
		})
	}
}
