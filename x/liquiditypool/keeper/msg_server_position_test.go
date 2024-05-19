package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestPositionMsgServerCreate(t *testing.T) {
	_, srv, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	sender := "A"
	for i := 0; i < 5; i++ {
		resp, err := srv.CreatePosition(wctx, &types.MsgCreatePosition{Sender: sender})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))
	}
}

func TestPositionMsgServerIncreaseLiquidity(t *testing.T) {
	sender := "A"

	tests := []struct {
		desc    string
		request *types.MsgIncreaseLiquidity
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgIncreaseLiquidity{Sender: sender},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgIncreaseLiquidity{Sender: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgIncreaseLiquidity{Sender: sender, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			_, err := srv.CreatePosition(wctx, &types.MsgCreatePosition{Sender: sender})
			require.NoError(t, err)

			_, err = srv.IncreaseLiquidity(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPositionMsgServerDecreaseLiquidity(t *testing.T) {
	sender := "A"

	tests := []struct {
		desc    string
		request *types.MsgDecreaseLiquidity
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDecreaseLiquidity{Sender: sender},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDecreaseLiquidity{Sender: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgDecreaseLiquidity{Sender: sender, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			_, err := srv.CreatePosition(wctx, &types.MsgCreatePosition{Sender: sender})
			require.NoError(t, err)
			_, err = srv.DecreaseLiquidity(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
