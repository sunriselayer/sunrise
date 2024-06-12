package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestPositionMsgServerCreate(t *testing.T) {
	k, bk, srv, ctx := setupMsgServer(t)
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

	// TODO: mint coins before position creation_
	_ = bk
	// err = bk.MintCoins(ctx, minttypes.ModuleName, sdk.Coins{sdk.NewInt64Coin("base", 10000000), sdk.NewInt64Coin("quote", 10000000)})
	// require.NoError(t, err)

	// Create positions
	for i := 0; i < 5; i++ {
		_, err := srv.CreatePosition(wctx, &types.MsgCreatePosition{
			Sender:         sender.String(),
			PoolId:         0,
			LowerTick:      0,
			UpperTick:      1,
			TokenBase:      sdk.NewInt64Coin("base", 10000),
			TokenQuote:     sdk.NewInt64Coin("quote", 10000),
			MinAmountBase:  math.NewInt(1),
			MinAmountQuote: math.NewInt(1),
		})
		require.Error(t, err)
		_ = k
		// require.Equal(t, i, int(resp.Id))
	}
}

func TestPositionMsgServerIncreaseLiquidity(t *testing.T) {
	sender := "sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f"

	tests := []struct {
		desc    string
		request *types.MsgIncreaseLiquidity
		err     error
	}{
		// {
		// 	desc: "Completed",
		// 	request: &types.MsgIncreaseLiquidity{
		// 		Sender:         sender,
		// 		Id:             0,
		// 		AmountBase:     math.NewInt(1),
		// 		AmountQuote:    math.NewInt(1),
		// 		MinAmountBase:  math.NewInt(1),
		// 		MinAmountQuote: math.NewInt(1),
		// 	},
		// },
		{
			desc: "Unauthorized",
			request: &types.MsgIncreaseLiquidity{
				Sender:         "B",
				Id:             0,
				AmountBase:     math.NewInt(1),
				AmountQuote:    math.NewInt(1),
				MinAmountBase:  math.NewInt(1),
				MinAmountQuote: math.NewInt(1),
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
				MinAmountBase:  math.NewInt(1),
				MinAmountQuote: math.NewInt(1),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, _, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			k.SetPool(wctx, types.Pool{
				Id:         0,
				DenomBase:  "base",
				DenomQuote: "quote",
				FeeRate:    math.LegacyNewDecWithPrec(1, 2),
				TickParams: types.TickParams{
					PriceRatio: math.LegacyNewDecWithPrec(10001, 4),
					BaseOffset: math.LegacyNewDecWithPrec(5, 1),
				},
				CurrentTick:          0,
				CurrentTickLiquidity: math.LegacyOneDec(),
				CurrentSqrtPrice:     math.LegacyOneDec(),
			})
			k.SetPosition(wctx, types.Position{
				Id:        0,
				Address:   sender,
				LowerTick: 0,
				UpperTick: 1,
				Liquidity: math.LegacyNewDec(10000),
			})
			_, err := srv.IncreaseLiquidity(wctx, tc.request)
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
		// {
		// 	desc: "Completed",
		// 	request: &types.MsgDecreaseLiquidity{
		// 		Sender:    sender,
		// 		Id:        0,
		// 		Liquidity: math.LegacyOneDec(),
		// 	},
		// },
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
			k, _, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			k.SetPool(wctx, types.Pool{
				Id:         0,
				DenomBase:  "base",
				DenomQuote: "quote",
				FeeRate:    math.LegacyNewDecWithPrec(1, 2),
				TickParams: types.TickParams{
					PriceRatio: math.LegacyNewDecWithPrec(10001, 4),
					BaseOffset: math.LegacyNewDecWithPrec(5, 1),
				},
				CurrentTick:          0,
				CurrentTickLiquidity: math.LegacyOneDec(),
				CurrentSqrtPrice:     math.LegacyOneDec(),
			})
			k.SetPosition(wctx, types.Position{
				Id:        0,
				Address:   sender,
				LowerTick: 0,
				UpperTick: 1,
				Liquidity: math.LegacyNewDec(10000),
			})
			_, err := srv.DecreaseLiquidity(wctx, tc.request)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
