package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestPositionMsgServerCreate(t *testing.T) {
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

	// Create positions
	for i := 0; i < 5; i++ {
		_, err := srv.CreatePosition(wctx, &types.MsgCreatePosition{
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
	}
}

func TestPositionMsgServerIncreaseLiquidity(t *testing.T) {
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
				AmountBase:     math.NewInt(1),
				AmountQuote:    math.NewInt(1),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgIncreaseLiquidity{
				Sender:         "B",
				Id:             0,
				AmountBase:     math.NewInt(1),
				AmountQuote:    math.NewInt(1),
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
				AmountBase:     math.NewInt(1),
				AmountQuote:    math.NewInt(1),
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
				LowerTick:      0,
				UpperTick:      1,
				TokenBase:      sdk.NewInt64Coin("base", 10000),
				TokenQuote:     sdk.NewInt64Coin("quote", 10000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
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
	sender := "sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f"

	tests := []struct {
		desc    string
		request *types.MsgDecreaseLiquidity
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDecreaseLiquidity{
				Sender:    sender,
				Id:        0,
				Liquidity: math.LegacyOneDec(),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDecreaseLiquidity{
				Sender:    "B",
				Id:        0,
				Liquidity: math.LegacyOneDec(),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDecreaseLiquidity{
				Sender:    sender,
				Id:        10,
				Liquidity: math.LegacyOneDec(),
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
				LowerTick:      0,
				UpperTick:      1,
				TokenBase:      sdk.NewInt64Coin("base", 10000),
				TokenQuote:     sdk.NewInt64Coin("quote", 10000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			_, err = srv.DecreaseLiquidity(wctx, tc.request)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
