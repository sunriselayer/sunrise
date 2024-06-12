package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, bankkeeper.Keeper, types.MsgServer, context.Context) {
	k, bk, ctx := keepertest.LiquiditypoolKeeper(t)
	return k, bk, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, _, ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}
