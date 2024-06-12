package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestPoolMsgServerCreate(t *testing.T) {
	k, srv, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	sender := sdk.AccAddress("sender")
	resp, err := srv.CreatePool(wctx, &types.MsgCreatePool{
		Authority:  sender.String(),
		DenomBase:  "base",
		DenomQuote: "quote",
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "0.5",
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.Id)

	// check created pool and status
	pool, found := k.GetPool(ctx, resp.Id)
	require.True(t, found)
	require.Equal(t, pool.TickParams.PriceRatio.String(), "1.000100000000000000")
	require.Equal(t, pool.TickParams.BaseOffset.String(), "0.500000000000000000")
	require.Equal(t, pool.FeeRate.String(), "0.010000000000000000")
	require.Equal(t, pool.CurrentTick, int64(0))
	require.Equal(t, pool.CurrentTickLiquidity.String(), "0.000000000000000000")
	require.Equal(t, pool.CurrentSqrtPrice.String(), "0.000000000000000000")
	require.Equal(t, pool.DenomBase, "base")
	require.Equal(t, pool.DenomQuote, "quote")

	// try creating another pool with same info
	_, err = srv.CreatePool(wctx, &types.MsgCreatePool{
		Authority:  sender.String(),
		DenomBase:  "base",
		DenomQuote: "quote",
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "0.5",
	})
	require.NoError(t, err)
}
