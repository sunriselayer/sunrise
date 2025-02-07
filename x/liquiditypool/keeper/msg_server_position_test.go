package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestMsgServerCreatePosition(t *testing.T) {
	_, bk, srv, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	sender := sdk.AccAddress("sender")
	_, err := srv.CreatePool(wctx, &types.MsgCreatePool{
		Authority:  sender.String(),
		DenomBase:  "base",
		DenomQuote: "quote",
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "0.5",
	})
	require.NoError(t, err)

	bk.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	bk.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Create 1st position
	resp, err := srv.CreatePosition(wctx, &types.MsgCreatePosition{
		Sender:         sender.String(),
		PoolId:         0,
		LowerTick:      0,
		UpperTick:      1,
		TokenBase:      sdk.NewInt64Coin("base", 10000),
		TokenQuote:     sdk.NewInt64Coin("quote", 10000),
		MinAmountBase:  math.NewInt(0),
		MinAmountQuote: math.NewInt(0),
	})
	require.NoError(t, err)
	require.Equal(t, resp.Id, uint64(0))
	require.Equal(t, resp.AmountBase.String(), "10001")
	require.Equal(t, resp.AmountQuote.String(), "0")
	require.Equal(t, resp.Liquidity.String(), "200020000.062502249619530703")

	// Create 2nd position with same tick
	resp, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
		Sender:         sender.String(),
		PoolId:         0,
		LowerTick:      0,
		UpperTick:      1,
		TokenBase:      sdk.NewInt64Coin("base", 10000),
		TokenQuote:     sdk.NewInt64Coin("quote", 10000),
		MinAmountBase:  math.NewInt(0),
		MinAmountQuote: math.NewInt(0),
	})
	require.NoError(t, err)
	require.Equal(t, resp.Id, uint64(1))
	require.Equal(t, resp.AmountBase.String(), "10001")
	require.Equal(t, resp.AmountQuote.String(), "0")
	require.Equal(t, resp.Liquidity.String(), "200020000.062502249619530703")

	// Create 3rd position with different tick
	resp, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
		Sender:         sender.String(),
		PoolId:         0,
		LowerTick:      -10,
		UpperTick:      10,
		TokenBase:      sdk.NewInt64Coin("base", 10000),
		TokenQuote:     sdk.NewInt64Coin("quote", 10000),
		MinAmountBase:  math.NewInt(0),
		MinAmountQuote: math.NewInt(0),
	})
	require.NoError(t, err)
	require.Equal(t, resp.Id, uint64(2))
	require.Equal(t, resp.AmountBase.String(), "10000")
	require.Equal(t, resp.AmountQuote.String(), "9048")
	require.Equal(t, resp.Liquidity.String(), "19053571.850177307210510444")
}

func TestMsgServerIncreaseLiquidity(t *testing.T) {
	sender := "sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f"

	tests := []struct {
		desc    string
		request *types.MsgIncreaseLiquidity
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgIncreaseLiquidity{
				Sender:         sender,
				Id:             0,
				AmountBase:     math.NewInt(100000),
				AmountQuote:    math.NewInt(100000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgIncreaseLiquidity{
				Sender:         "B",
				Id:             0,
				AmountBase:     math.NewInt(100000),
				AmountQuote:    math.NewInt(100000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "Unauthorized",
			request: &types.MsgIncreaseLiquidity{
				Sender:         sender,
				Id:             10,
				AmountBase:     math.NewInt(100000),
				AmountQuote:    math.NewInt(100000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, bk, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			bk.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			bk.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			_, err := srv.CreatePool(wctx, &types.MsgCreatePool{
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
				PoolId:         0,
				LowerTick:      -10,
				UpperTick:      10,
				TokenBase:      sdk.NewInt64Coin("base", 10000),
				TokenQuote:     sdk.NewInt64Coin("quote", 10000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			resp, err := srv.IncreaseLiquidity(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, resp.AmountBase.String(), "110000")
				require.Equal(t, resp.AmountQuote.String(), "99527")
				require.Equal(t, resp.PositionId, uint64(1))
			}
		})
	}
}

func TestMsgServerDecreaseLiquidity(t *testing.T) {
	sender := "sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f"
	_, bk, srv, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	bk.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	bk.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	_, err := srv.CreatePool(wctx, &types.MsgCreatePool{
		Authority:  sender,
		DenomBase:  "base",
		DenomQuote: "quote",
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "0.5",
	})
	require.NoError(t, err)

	resp, err := srv.CreatePosition(wctx, &types.MsgCreatePosition{
		Sender:         sender,
		PoolId:         0,
		LowerTick:      -10,
		UpperTick:      10,
		TokenBase:      sdk.NewInt64Coin("base", 10000),
		TokenQuote:     sdk.NewInt64Coin("quote", 10000),
		MinAmountBase:  math.NewInt(0),
		MinAmountQuote: math.NewInt(0),
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDecreaseLiquidity
		err     error
	}{
		{
			desc: "Successful deduction",
			request: &types.MsgDecreaseLiquidity{
				Sender:    sender,
				Id:        0,
				Liquidity: resp.Liquidity.String(),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDecreaseLiquidity{
				Sender:    "B",
				Id:        0,
				Liquidity: resp.Liquidity.String(),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "Position not found",
			request: &types.MsgDecreaseLiquidity{
				Sender:    sender,
				Id:        10,
				Liquidity: resp.Liquidity.String(),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			res, err := srv.DecreaseLiquidity(wctx, tc.request)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, res.AmountBase.String(), "10000")
				require.Equal(t, res.AmountQuote.String(), "9047")
			}
		})
	}
}
