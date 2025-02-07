package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/swap/keeper"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func setupMsgServer(t *testing.T) (keeper.Keeper, types.MsgServer, context.Context) {
	f := initFixture(t)
	ctx := sdk.UnwrapSDKContext(f.ctx)
	k := f.keeper
	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}
