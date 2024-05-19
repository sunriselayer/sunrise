package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestPoolMsgServerCreate(t *testing.T) {
	_, srv, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	sender := "A"
	for i := 0; i < 5; i++ {
		resp, err := srv.CreatePool(wctx, &types.MsgCreatePool{Authority: sender})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))
	}
}
