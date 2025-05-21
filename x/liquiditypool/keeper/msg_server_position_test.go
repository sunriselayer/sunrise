package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestMsgServerCreatePosition(t *testing.T) {
	_, mocks, srv, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).Return(consts.FeeDenom, nil).AnyTimes()

	sender := sdk.AccAddress("sender")
	quoteDenom := consts.FeeDenom
	_, err := srv.CreatePool(wctx, &types.MsgCreatePool{
		Sender:     sender.String(),
		DenomBase:  "base",
		DenomQuote: quoteDenom,
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "-0.5",
	})
	require.NoError(t, err)

	mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Create 1st position
	resp, err := srv.CreatePosition(wctx, &types.MsgCreatePosition{
		Sender:         sender.String(),
		PoolId:         0,
		LowerTick:      0,
		UpperTick:      1,
		TokenBase:      sdk.NewInt64Coin("base", 10000),
		TokenQuote:     sdk.NewInt64Coin(quoteDenom, 10000),
		MinAmountBase:  math.NewInt(0),
		MinAmountQuote: math.NewInt(0),
	})
	require.NoError(t, err)
	require.Equal(t, resp.Id, uint64(0))
	require.Equal(t, resp.AmountBase.String(), "10000")
	require.Equal(t, resp.AmountQuote.String(), "10000")
	require.Equal(t, resp.Liquidity, "400024997.187355990380704326")

	// Create 2nd position with same tick
	resp, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
		Sender:         sender.String(),
		PoolId:         0,
		LowerTick:      0,
		UpperTick:      1,
		TokenBase:      sdk.NewInt64Coin("base", 10000),
		TokenQuote:     sdk.NewInt64Coin(quoteDenom, 10000),
		MinAmountBase:  math.NewInt(0),
		MinAmountQuote: math.NewInt(0),
	})
	require.NoError(t, err)
	require.Equal(t, resp.Id, uint64(1))
	require.Equal(t, resp.AmountBase.String(), "10000")
	require.Equal(t, resp.AmountQuote.String(), "10000")
	require.Equal(t, resp.Liquidity, "400024997.187355990380704326")

	// Create 3rd position with different tick
	resp, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
		Sender:         sender.String(),
		PoolId:         0,
		LowerTick:      -10,
		UpperTick:      10,
		TokenBase:      sdk.NewInt64Coin("base", 10000),
		TokenQuote:     sdk.NewInt64Coin(quoteDenom, 10000),
		MinAmountBase:  math.NewInt(0),
		MinAmountQuote: math.NewInt(0),
	})
	require.NoError(t, err)
	require.Equal(t, resp.Id, uint64(2))
	require.Equal(t, resp.AmountBase.String(), "9048")
	require.Equal(t, resp.AmountQuote.String(), "10000")
	require.Equal(t, resp.Liquidity, "19053571.855846596797818151")
}

func TestMsgServerIncreaseLiquidity(t *testing.T) {
	initFixture(t)
	sender := sdk.AccAddress("sender")
	quoteDenom := consts.FeeDenom
	tests := []struct {
		desc    string
		request *types.MsgIncreaseLiquidity
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgIncreaseLiquidity{
				Sender:         sender.String(),
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
				Sender:         sdk.AccAddress("hoge").String(),
				Id:             0,
				AmountBase:     math.NewInt(100000),
				AmountQuote:    math.NewInt(100000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "Invalid ID",
			request: &types.MsgIncreaseLiquidity{
				Sender:         sender.String(),
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
			_, mocks, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).Return(consts.FeeDenom, nil).AnyTimes()
			_, err := srv.CreatePool(wctx, &types.MsgCreatePool{
				Sender:     sender.String(),
				DenomBase:  "base",
				DenomQuote: quoteDenom,
				FeeRate:    "0.01",
				PriceRatio: "1.0001",
				BaseOffset: "-0.5",
			})
			require.NoError(t, err)

			_, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
				Sender:         sender.String(),
				PoolId:         0,
				LowerTick:      -4155,
				UpperTick:      4054,
				TokenBase:      sdk.NewInt64Coin("base", 10000),
				TokenQuote:     sdk.NewInt64Coin(quoteDenom, 10000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			resp, err := srv.IncreaseLiquidity(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, resp.AmountBase.String(), "106400")
				require.Equal(t, resp.AmountQuote.String(), "110000")
				require.Equal(t, resp.PositionId, uint64(1))
			}
		})
	}
}

func TestMsgServerDecreaseLiquidity(t *testing.T) {
	initFixture(t)
	sender := sdk.AccAddress("sender")
	quoteDenom := consts.FeeDenom
	_, mocks, srv, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).Return(consts.FeeDenom, nil).AnyTimes()
	_, err := srv.CreatePool(wctx, &types.MsgCreatePool{
		Sender:     sender.String(),
		DenomBase:  "base",
		DenomQuote: quoteDenom,
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "-0.5",
	})
	require.NoError(t, err)

	resp, err := srv.CreatePosition(wctx, &types.MsgCreatePosition{
		Sender:         sender.String(),
		PoolId:         0,
		LowerTick:      -10,
		UpperTick:      10,
		TokenBase:      sdk.NewInt64Coin("base", 10000),
		TokenQuote:     sdk.NewInt64Coin(quoteDenom, 10000),
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
				Sender:    sender.String(),
				Id:        0,
				Liquidity: resp.Liquidity,
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDecreaseLiquidity{
				Sender:    "B",
				Id:        0,
				Liquidity: resp.Liquidity,
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "Position not found",
			request: &types.MsgDecreaseLiquidity{
				Sender:    sender.String(),
				Id:        10,
				Liquidity: resp.Liquidity,
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
				require.Equal(t, res.AmountBase.String(), "9047")
				require.Equal(t, res.AmountQuote.String(), "10000")
			}
		})
	}
}
