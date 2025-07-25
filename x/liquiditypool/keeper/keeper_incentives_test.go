package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestAllocateIncentive(t *testing.T) {
	sender := sdk.AccAddress("sender")
	quoteDenom := consts.StableDenom
	tests := []struct {
		desc            string
		poolId          uint64
		allocation      sdk.Coins
		expAccumulation string
		err             error
	}{
		{
			desc:            "Single token allocation",
			poolId:          0,
			allocation:      sdk.Coins{sdk.NewInt64Coin("xyz", 1000)},
			expAccumulation: "0.000052483597698410xyz",
		},
		{
			desc:            "Multiple tokens allocation",
			poolId:          0,
			allocation:      sdk.Coins{sdk.NewInt64Coin("uvw", 1000), sdk.NewInt64Coin("xyz", 1000)},
			expAccumulation: "0.000052483597698410uvw,0.000052483597698410xyz",
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
			poolId:          1,
			allocation:      sdk.Coins{sdk.NewInt64Coin("xyz", 1000)},
			expAccumulation: "",
			err:             types.ErrEmptyLiquidity,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, mocks, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			// First pool
			_, err := srv.CreatePool(wctx, &types.MsgCreatePool{
				Sender:     sender.String(),
				DenomBase:  "base",
				DenomQuote: quoteDenom,
				FeeRate:    "0.01",
				PriceRatio: "1.0001",
				BaseOffset: "-0.5",
			})
			require.NoError(t, err)

			// Second pool
			_, err = srv.CreatePool(wctx, &types.MsgCreatePool{
				Sender:     sender.String(),
				DenomBase:  "base",
				DenomQuote: quoteDenom,
				FeeRate:    "0.01",
				PriceRatio: "1.0001",
				BaseOffset: "-0.5",
			})
			require.NoError(t, err)

			_, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
				Sender:         sender.String(),
				PoolId:         0,
				LowerTick:      -10,
				UpperTick:      10,
				TokenBase:      sdk.NewInt64Coin("base", 10000),
				TokenQuote:     sdk.NewInt64Coin(quoteDenom, 10000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			err = k.AllocateIncentive(wctx, tc.poolId, sender, tc.allocation)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)

				resp, err := k.GetFeeAccumulator(wctx, tc.poolId)
				require.NoError(t, err)
				require.Equal(t, resp.AccumValue.String(), tc.expAccumulation)
				require.Equal(t, resp.TotalShares, "19053571.855846596797818151")
			}
		})
	}
}
