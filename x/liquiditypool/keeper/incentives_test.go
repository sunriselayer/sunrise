package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestAllocateIncentive(t *testing.T) {
	sender := "sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f"
	senderAcc := sdk.MustAccAddressFromBech32(sender)

	tests := []struct {
		desc            string
		poolId          uint64
		allocation      sdk.Coins
		expAccumulation string
		err             error
	}{
		{
			desc:            "Single token allocation",
			poolId:          1,
			allocation:      sdk.Coins{sdk.NewInt64Coin("xyz", 1000)},
			expAccumulation: "0.000052483597714026xyz",
		},
		{
			desc:            "Multiple tokens allocation",
			poolId:          1,
			allocation:      sdk.Coins{sdk.NewInt64Coin("uvw", 1000), sdk.NewInt64Coin("xyz", 1000)},
			expAccumulation: "0.000052483597714026uvw,0.000052483597714026xyz",
		},
		{
			desc:            "Not available pool",
			poolId:          3,
			allocation:      sdk.Coins{sdk.NewInt64Coin("xyz", 1000)},
			expAccumulation: "",
			err:             types.ErrPoolNotFound,
		},
		{
			desc:            "Not available position",
			poolId:          2,
			allocation:      sdk.Coins{sdk.NewInt64Coin("xyz", 1000)},
			expAccumulation: "",
			err:             types.ErrEmptyLiquidity,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, bk, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			bk.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			bk.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			// First pool
			_, err := srv.CreatePool(wctx, &types.MsgCreatePool{
				Authority:  sender,
				DenomBase:  "base",
				DenomQuote: "quote",
				FeeRate:    "0.01",
				PriceRatio: "1.0001",
				BaseOffset: "0.5",
			})
			require.NoError(t, err)

			// Second pool
			_, err = srv.CreatePool(wctx, &types.MsgCreatePool{
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

			err = k.AllocateIncentive(wctx, tc.poolId, senderAcc, tc.allocation)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)

				resp, err := k.GetFeeAccumulator(wctx, tc.poolId)
				require.NoError(t, err)
				require.Equal(t, resp.AccumValue.String(), tc.expAccumulation)
				require.Equal(t, resp.TotalShares.String(), "19053571.850177307210510444")
			}
		})
	}
}
