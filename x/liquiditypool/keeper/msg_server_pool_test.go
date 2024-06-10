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

	sender := sdk.AccAddress("sender")
	for i := 0; i < 1; i++ {
		resp, err := srv.CreatePool(wctx, &types.MsgCreatePool{
			Authority:  sender.String(),
			DenomBase:  "base",
			DenomQuote: "quote",
			FeeRate:    "0.01",
			PriceRatio: "1.0001",
			BaseOffset: "0.5",
		})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))
	}
}
