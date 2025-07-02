package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestPoolMsgServerCreate(t *testing.T) {
	k, mocks, srv, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).Return(consts.StableDenom, nil).AnyTimes()

	sender := sdk.AccAddress("sender")
	quoteDenom := consts.StableDenom

	resp, err := srv.CreatePool(wctx, &types.MsgCreatePool{
		Sender:     sender.String(),
		DenomBase:  "base",
		DenomQuote: quoteDenom,
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "-0.5",
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.Id)

	// check created pool and status
	pool, found, err := k.GetPool(ctx, resp.Id)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, pool.TickParams.PriceRatio, "1.000100000000000000")
	require.Equal(t, pool.TickParams.BaseOffset, "-0.500000000000000000")
	require.Equal(t, pool.FeeRate, "0.010000000000000000")
	require.Equal(t, pool.CurrentTick, int64(0))
	require.Equal(t, pool.CurrentTickLiquidity, "0.000000000000000000")
	require.Equal(t, pool.CurrentSqrtPrice, "0.000000000000000000")
	require.Equal(t, pool.DenomBase, "base")
	require.Equal(t, pool.DenomQuote, quoteDenom)

	// try creating another pool with same info
	_, err = srv.CreatePool(wctx, &types.MsgCreatePool{
		Sender:     sender.String(),
		DenomBase:  "base",
		DenomQuote: quoteDenom,
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "-0.5",
	})
	require.NoError(t, err)
}
