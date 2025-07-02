package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// TODO: add test for updateFeeGrowthGlobal
// TODO: add test for swapOutAmtGivenIn
// TODO: add test for swapInAmtGivenOut
// TODO: add test for iteratorToNextTickSqrtPriceTarget
// TODO: add test for computeOutAmtGivenIn
// TODO: add test for computeInAmtGivenOut
// TODO: add test for swapCrossTickLogic
// TODO: add test for updatePoolForSwap
// TODO: add test for setupSwapHelper
// TODO: add test for validateSwapProgressAndAmountConsumption
// TODO: add test for edgeCaseInequalityBasedOnSwapHelper
// TODO: add test for ComputeMaxInAmtGivenMaxTicksCrossed

func TestSwapExactAmountIn_SinglePosition(t *testing.T) {
	sender := sdk.AccAddress("sender")
	quoteDenom := consts.StableDenom

	tests := []struct {
		desc         string
		tokenIn      sdk.Coin
		denomOut     string
		feeEnabled   bool
		expAmountOut math.Int
		expTickIndex int64
		err          error
	}{
		{
			desc:         "Base to quote",
			tokenIn:      sdk.NewInt64Coin("base", 100000),
			denomOut:     quoteDenom,
			feeEnabled:   false,
			expAmountOut: math.NewInt(99994),
			expTickIndex: -1,
		},
		{
			desc:         "Quote to base",
			tokenIn:      sdk.NewInt64Coin("base", 100000),
			denomOut:     quoteDenom,
			feeEnabled:   false,
			expAmountOut: math.NewInt(99994),
			expTickIndex: -1,
		},
		{
			desc:         "Fee enabled",
			tokenIn:      sdk.NewInt64Coin("base", 100000),
			denomOut:     quoteDenom,
			feeEnabled:   true,
			expAmountOut: math.NewInt(98994),
			expTickIndex: -1,
		},
		{
			desc:       "Ran out of ticks",
			tokenIn:    sdk.NewInt64Coin("base", 1000000000),
			denomOut:   quoteDenom,
			feeEnabled: true,
			err:        types.ErrRanOutOfTicks,
		},
		{
			desc:         "Empty token in",
			tokenIn:      sdk.NewInt64Coin("base", 0),
			denomOut:     quoteDenom,
			feeEnabled:   true,
			expAmountOut: math.NewInt(100000),
			err:          types.ErrUnexpectedCalcAmount,
		},
		{
			desc:         "same token in and out",
			tokenIn:      sdk.NewInt64Coin("base", 0),
			denomOut:     "base",
			feeEnabled:   true,
			expAmountOut: math.NewInt(100000),
			err:          types.ErrDenomDuplication,
		},
		{
			desc:         "invalid denomOut",
			tokenIn:      sdk.NewInt64Coin("base", 0),
			denomOut:     "invalid",
			feeEnabled:   true,
			expAmountOut: math.NewInt(100000),
			err:          types.ErrInvalidOutDenom,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, mocks, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).Return(consts.StableDenom, nil).AnyTimes()

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
				LowerTick:      -10,
				UpperTick:      10,
				TokenBase:      sdk.NewInt64Coin("base", 1000000),
				TokenQuote:     sdk.NewInt64Coin(quoteDenom, 1000000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			pool, found, err := k.GetPool(ctx, 0)
			require.NoError(t, err)
			require.True(t, found)

			amountOut, err := k.SwapExactAmountIn(wctx, sender, pool, tc.tokenIn, tc.denomOut, tc.feeEnabled)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)

				pool, found, err = k.GetPool(ctx, 0)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, amountOut.String(), tc.expAmountOut.String())
				require.Equal(t, pool.CurrentTick, tc.expTickIndex)
			}
		})
	}
}

func TestSwapExactAmountIn_MultiplePositions(t *testing.T) {
	sender := sdk.AccAddress("sender")
	quoteDenom := consts.StableDenom

	tests := []struct {
		desc         string
		tokenIn      sdk.Coin
		denomOut     string
		feeEnabled   bool
		expAmountOut math.Int
		expTickIndex int64
		err          error
	}{
		{
			desc:         "Swap on multiple positions pool",
			tokenIn:      sdk.NewInt64Coin("base", 110000),
			denomOut:     quoteDenom,
			feeEnabled:   false,
			expAmountOut: math.NewInt(109947),
			expTickIndex: -10,
		},
		{
			desc:       "Ran out of ticks",
			tokenIn:    sdk.NewInt64Coin("base", 1000000000),
			denomOut:   quoteDenom,
			feeEnabled: true,
			err:        types.ErrRanOutOfTicks,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, mocks, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).Return(consts.StableDenom, nil).AnyTimes()

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
				LowerTick:      -10,
				UpperTick:      10,
				TokenBase:      sdk.NewInt64Coin("base", 100000),
				TokenQuote:     sdk.NewInt64Coin(quoteDenom, 100000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			_, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
				Sender:         sender.String(),
				PoolId:         0,
				LowerTick:      -50,
				UpperTick:      50,
				TokenBase:      sdk.NewInt64Coin("base", 100000),
				TokenQuote:     sdk.NewInt64Coin(quoteDenom, 100000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			pool, found, err := k.GetPool(ctx, 0)
			require.NoError(t, err)
			require.True(t, found)

			amountOut, err := k.SwapExactAmountIn(wctx, sender, pool, tc.tokenIn, tc.denomOut, tc.feeEnabled)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, amountOut.String(), tc.expAmountOut.String())

				pool, found, err = k.GetPool(ctx, 0)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, pool.CurrentTick, tc.expTickIndex)
			}
		})
	}
}

func TestSwapExactAmountOut_SinglePosition(t *testing.T) {
	sender := sdk.AccAddress("sender")
	quoteDenom := consts.StableDenom

	tests := []struct {
		desc         string
		tokenOut     sdk.Coin
		denomIn      string
		feeEnabled   bool
		expAmountIn  math.Int
		expTickIndex int64
		err          error
	}{
		{
			desc:         "Base to quote",
			tokenOut:     sdk.NewInt64Coin("base", 100000),
			denomIn:      quoteDenom,
			feeEnabled:   false,
			expAmountIn:  math.NewInt(100006),
			expTickIndex: 1,
		},
		{
			desc:         "Quote to base",
			tokenOut:     sdk.NewInt64Coin("base", 100000),
			denomIn:      quoteDenom,
			feeEnabled:   false,
			expAmountIn:  math.NewInt(100006),
			expTickIndex: 1,
		},
		{
			desc:         "Fee enabled",
			tokenOut:     sdk.NewInt64Coin("base", 100000),
			denomIn:      quoteDenom,
			feeEnabled:   true,
			expAmountIn:  math.NewInt(101017),
			expTickIndex: 1,
		},
		{
			desc:       "Ran out of ticks",
			tokenOut:   sdk.NewInt64Coin("base", 1000000000),
			denomIn:    quoteDenom,
			feeEnabled: true,
			err:        types.ErrRanOutOfTicks,
		},
		{
			desc:        "Empty token in",
			tokenOut:    sdk.NewInt64Coin("base", 0),
			denomIn:     quoteDenom,
			feeEnabled:  true,
			expAmountIn: math.NewInt(100000),
			err:         types.ErrUnexpectedCalcAmount,
		},
		{
			desc:        "same token in and out",
			tokenOut:    sdk.NewInt64Coin("base", 0),
			denomIn:     "base",
			feeEnabled:  true,
			expAmountIn: math.NewInt(100000),
			err:         types.ErrDenomDuplication,
		},
		{
			desc:        "invalid denomIn",
			tokenOut:    sdk.NewInt64Coin("base", 0),
			denomIn:     "invalid",
			feeEnabled:  true,
			expAmountIn: math.NewInt(100000),
			err:         types.ErrInvalidInDenom,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, mocks, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).Return(consts.StableDenom, nil).AnyTimes()

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
				LowerTick:      -10,
				UpperTick:      10,
				TokenBase:      sdk.NewInt64Coin("base", 1000000),
				TokenQuote:     sdk.NewInt64Coin(quoteDenom, 1000000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			pool, found, err := k.GetPool(ctx, 0)
			require.NoError(t, err)
			require.True(t, found)

			amountIn, err := k.SwapExactAmountOut(wctx, sender, pool, tc.tokenOut, tc.denomIn, tc.feeEnabled)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)

				pool, found, err = k.GetPool(ctx, 0)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, amountIn.String(), tc.expAmountIn.String())
				require.Equal(t, pool.CurrentTick, tc.expTickIndex)
			}
		})
	}
}

func TestSwapExactAmountOut_MultiplePositions(t *testing.T) {
	sender := sdk.AccAddress("sender")
	quoteDenom := consts.StableDenom

	tests := []struct {
		desc         string
		tokenOut     sdk.Coin
		denomIn      string
		feeEnabled   bool
		expAmountIn  math.Int
		expTickIndex int64
		err          error
	}{
		{
			desc:         "Swap on multiple positions pool",
			tokenOut:     sdk.NewInt64Coin("base", 110000),
			denomIn:      quoteDenom,
			feeEnabled:   false,
			expAmountIn:  math.NewInt(110054),
			expTickIndex: 10,
		},
		{
			desc:       "Ran out of ticks",
			tokenOut:   sdk.NewInt64Coin("base", 1000000000),
			denomIn:    quoteDenom,
			feeEnabled: true,
			err:        types.ErrRanOutOfTicks,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, mocks, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).Return(consts.StableDenom, nil).AnyTimes()

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
				LowerTick:      -10,
				UpperTick:      10,
				TokenBase:      sdk.NewInt64Coin("base", 100000),
				TokenQuote:     sdk.NewInt64Coin(quoteDenom, 100000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			_, err = srv.CreatePosition(wctx, &types.MsgCreatePosition{
				Sender:         sender.String(),
				PoolId:         0,
				LowerTick:      -50,
				UpperTick:      50,
				TokenBase:      sdk.NewInt64Coin("base", 100000),
				TokenQuote:     sdk.NewInt64Coin(quoteDenom, 100000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			pool, found, err := k.GetPool(ctx, 0)
			require.NoError(t, err)
			require.True(t, found)

			amountIn, err := k.SwapExactAmountOut(wctx, sender, pool, tc.tokenOut, tc.denomIn, tc.feeEnabled)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, amountIn.String(), tc.expAmountIn.String())

				pool, found, err = k.GetPool(ctx, 0)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, pool.CurrentTick, tc.expTickIndex)
			}
		})
	}
}

func TestGetValidatedPoolAndAccumulator(t *testing.T) {
	f := initFixture(t)
	ctx := sdk.UnwrapSDKContext(f.ctx)
	k := f.keeper
	quoteDenom := consts.StableDenom
	// when pool does not exist
	_, _, err := k.GetValidatedPoolAndAccumulator(ctx, 1, "base", quoteDenom)
	require.Error(t, err)

	// when accumulator does not exist
	err = k.SetPool(ctx, types.Pool{
		Id:                   1,
		DenomBase:            "base",
		DenomQuote:           quoteDenom,
		FeeRate:              math.LegacyZeroDec().String(),
		TickParams:           types.TickParams{},
		CurrentTick:          0,
		CurrentTickLiquidity: math.LegacyOneDec().String(),
		CurrentSqrtPrice:     math.LegacyOneDec().String(),
	})
	if err != nil {
		t.Fatalf("failed to set pool: %v", err)
	}
	_, _, err = k.GetValidatedPoolAndAccumulator(ctx, 1, "base", quoteDenom)
	require.Error(t, err)

	// When both accumulator and pool exist
	err = k.InitAccumulator(ctx, "fee_pool_accumulator/1")
	require.NoError(t, err)
	pool, accumulator, err := k.GetValidatedPoolAndAccumulator(ctx, 1, "base", quoteDenom)
	require.NoError(t, err)
	require.Equal(t, pool.Id, uint64(1))
	require.Equal(t, pool.DenomBase, "base")
	require.Equal(t, pool.DenomQuote, quoteDenom)
	require.Equal(t, pool.FeeRate, "0.000000000000000000")
	require.Equal(t, pool.CurrentTick, int64(0))
	require.Equal(t, pool.CurrentTickLiquidity, "1.000000000000000000")
	require.Equal(t, pool.CurrentSqrtPrice, "1.000000000000000000")
	require.Equal(t, accumulator.AccumValue.String(), "")
	require.Equal(t, accumulator.Name, "fee_pool_accumulator/1")
	require.Equal(t, accumulator.TotalShares, "0.000000000000000000")

	// When invalid denom's put
	_, _, err = k.GetValidatedPoolAndAccumulator(ctx, 1, "invalid_denom", quoteDenom)
	require.Error(t, err)
}

func TestCalculateResultExactAmountOut(t *testing.T) {
	sender := sdk.AccAddress("sender")
	quoteDenom := consts.StableDenom

	tests := []struct {
		desc         string
		tokenOut     sdk.Coin
		denomIn      string
		feeEnabled   bool
		expAmountIn  math.Int
		expTickIndex int64
		err          error
	}{
		{
			desc:         "Base to quote",
			tokenOut:     sdk.NewInt64Coin("base", 100000),
			denomIn:      quoteDenom,
			feeEnabled:   false,
			expAmountIn:  math.NewInt(100006),
			expTickIndex: 0,
		},
		{
			desc:         "Quote to base",
			tokenOut:     sdk.NewInt64Coin("base", 100000),
			denomIn:      quoteDenom,
			feeEnabled:   false,
			expAmountIn:  math.NewInt(100006),
			expTickIndex: 0,
		},
		{
			desc:         "Fee enabled",
			tokenOut:     sdk.NewInt64Coin("base", 100000),
			denomIn:      quoteDenom,
			feeEnabled:   true,
			expAmountIn:  math.NewInt(101017),
			expTickIndex: 0,
		},
		{
			desc:       "Ran out of ticks",
			tokenOut:   sdk.NewInt64Coin("base", 1000000000),
			denomIn:    quoteDenom,
			feeEnabled: true,
			err:        types.ErrRanOutOfTicks,
		},
		{
			desc:        "Empty token in",
			tokenOut:    sdk.NewInt64Coin("base", 0),
			denomIn:     quoteDenom,
			feeEnabled:  true,
			expAmountIn: math.NewInt(0),
			err:         nil,
		},
		{
			desc:        "same token in and out",
			tokenOut:    sdk.NewInt64Coin("base", 0),
			denomIn:     "base",
			feeEnabled:  true,
			expAmountIn: math.NewInt(100000),
			err:         types.ErrDenomDuplication,
		},
		{
			desc:        "invalid denomIn",
			tokenOut:    sdk.NewInt64Coin("base", 0),
			denomIn:     "invalid",
			feeEnabled:  true,
			expAmountIn: math.NewInt(100000),
			err:         types.ErrInvalidInDenom,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, mocks, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).Return(consts.StableDenom, nil).AnyTimes()

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
				LowerTick:      -10,
				UpperTick:      10,
				TokenBase:      sdk.NewInt64Coin("base", 1000000),
				TokenQuote:     sdk.NewInt64Coin(quoteDenom, 1000000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			pool, found, err := k.GetPool(ctx, 0)
			require.NoError(t, err)
			require.True(t, found)

			amountIn, err := k.CalculateResultExactAmountOut(wctx, pool, tc.tokenOut, tc.denomIn, tc.feeEnabled)
			if tc.err != nil {
				require.Error(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, amountIn.String(), tc.expAmountIn.String())
			}
		})
	}
}

func TestCalculateResultExactAmountIn(t *testing.T) {
	sender := sdk.AccAddress("sender")
	quoteDenom := consts.StableDenom

	tests := []struct {
		desc         string
		tokenIn      sdk.Coin
		denomOut     string
		feeEnabled   bool
		expAmountOut math.Int
		expTickIndex int64
		err          error
	}{
		{
			desc:         "Base to quote",
			tokenIn:      sdk.NewInt64Coin("base", 100000),
			denomOut:     quoteDenom,
			feeEnabled:   false,
			expAmountOut: math.NewInt(99994),
			expTickIndex: -2,
		},
		{
			desc:         "Quote to base",
			tokenIn:      sdk.NewInt64Coin("base", 100000),
			denomOut:     quoteDenom,
			feeEnabled:   false,
			expAmountOut: math.NewInt(99994),
			expTickIndex: -2,
		},
		{
			desc:         "Fee enabled",
			tokenIn:      sdk.NewInt64Coin("base", 100000),
			denomOut:     quoteDenom,
			feeEnabled:   true,
			expAmountOut: math.NewInt(98994),
			expTickIndex: -2,
		},
		{
			desc:       "Ran out of ticks",
			tokenIn:    sdk.NewInt64Coin("base", 1000000000),
			denomOut:   quoteDenom,
			feeEnabled: true,
			err:        types.ErrRanOutOfTicks,
		},
		{
			desc:         "Empty token in",
			tokenIn:      sdk.NewInt64Coin("base", 0),
			denomOut:     quoteDenom,
			feeEnabled:   true,
			expAmountOut: math.NewInt(0),
			err:          nil,
		},
		{
			desc:         "same token in and out",
			tokenIn:      sdk.NewInt64Coin("base", 0),
			denomOut:     "base",
			feeEnabled:   true,
			expAmountOut: math.NewInt(100000),
			err:          types.ErrDenomDuplication,
		},
		{
			desc:         "invalid denomOut",
			tokenIn:      sdk.NewInt64Coin("base", 0),
			denomOut:     "invalid",
			feeEnabled:   true,
			expAmountOut: math.NewInt(100000),
			err:          types.ErrInvalidOutDenom,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, mocks, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).Return(consts.StableDenom, nil).AnyTimes()

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
				LowerTick:      -10,
				UpperTick:      10,
				TokenBase:      sdk.NewInt64Coin("base", 1000000),
				TokenQuote:     sdk.NewInt64Coin(quoteDenom, 1000000),
				MinAmountBase:  math.NewInt(0),
				MinAmountQuote: math.NewInt(0),
			})
			require.NoError(t, err)

			pool, found, err := k.GetPool(ctx, 0)
			require.NoError(t, err)
			require.True(t, found)

			amountOut, err := k.CalculateResultExactAmountIn(wctx, pool, tc.tokenIn, tc.denomOut, tc.feeEnabled)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, amountOut.String(), tc.expAmountOut.String())
			}
		})
	}
}
