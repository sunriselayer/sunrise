package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/testutil"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func setupMsgServer(t *testing.T) (keeper.Keeper, *testutil.MockBankKeeper, types.MsgServer, context.Context) {
	f := initFixture(t)
	ctx := sdk.UnwrapSDKContext(f.ctx)
	k := f.keeper
	bk := testutil.NewMockBankKeeper(gomock.NewController(t))
	return k, bk, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, _, ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}
