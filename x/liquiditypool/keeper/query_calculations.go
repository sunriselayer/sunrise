package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) CalculationCreatePosition(ctx context.Context, req *types.QueryCalculationCreatePositionRequest) (*types.QueryCalculationCreatePositionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	pool, found, err := q.k.GetPool(ctx, req.PoolId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, types.ErrPoolNotFound
	}

	lowerTick, ok := sdkmath.NewIntFromString(req.LowerTick)
	if !ok {
		return nil, types.ErrInvalidTickers
	}
	upperTick, ok := sdkmath.NewIntFromString(req.UpperTick)
	if !ok {
		return nil, types.ErrInvalidTickers
	}
	err = types.CheckTicks(lowerTick.Int64(), upperTick.Int64())
	if err != nil {
		return nil, types.ErrInvalidTickers
	}
	amount, ok := sdkmath.NewIntFromString(req.Amount)
	if !ok {
		return nil, types.ErrInvalidTokenAmounts
	}

	sqrtPriceLowerTick, sqrtPriceUpperTick, err := types.TicksToSqrtPrice(lowerTick.Int64(), upperTick.Int64(), pool.TickParams)
	if err != nil {
		return nil, err
	}

	currentSqrtPrice, err := math.LegacyNewDecFromStr(pool.CurrentSqrtPrice)
	if err != nil {
		return nil, err
	}
	var liquidityDelta math.LegacyDec
	if req.Denom == pool.DenomBase {
		liquidityDelta = types.LiquidityBase(amount, currentSqrtPrice, sqrtPriceUpperTick)
		_, actualAmountQuote, err := pool.CalcActualAmounts(lowerTick.Int64(), upperTick.Int64(), liquidityDelta)
		if err != nil {
			return nil, err
		}
		return &types.QueryCalculationCreatePositionResponse{
			Amount: sdk.NewCoin(pool.DenomQuote, actualAmountQuote.TruncateInt()),
		}, nil
	} else if req.Denom == pool.DenomQuote {
		liquidityDelta = types.LiquidityQuote(amount, currentSqrtPrice, sqrtPriceLowerTick)
		actualAmountBase, _, err := pool.CalcActualAmounts(lowerTick.Int64(), upperTick.Int64(), liquidityDelta)
		if err != nil {
			return nil, err
		}
		return &types.QueryCalculationCreatePositionResponse{
			Amount: sdk.NewCoin(pool.DenomBase, actualAmountBase.TruncateInt()),
		}, nil
	} else {
		return nil, types.ErrInvalidTickers
	}
}

func (q queryServer) CalculationIncreaseLiquidity(ctx context.Context, req *types.QueryCalculationIncreaseLiquidityRequest) (*types.QueryCalculationIncreaseLiquidityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	amountIn, ok := sdkmath.NewIntFromString(req.AmountIn)
	if !ok {
		return nil, types.ErrInvalidTokenAmounts
	}

	if amountIn.IsNegative() {
		return nil, types.ErrNegativeTokenAmount
	}

	if amountIn.IsZero() {
		return nil, types.ErrInvalidTokenAmounts
	}

	position, found, err := q.k.GetPosition(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, types.ErrPositionNotFound
	}
	pool, found, err := q.k.GetPool(ctx, position.PoolId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, types.ErrPoolNotFound
	}
	liquidity, err := math.LegacyNewDecFromStr(position.Liquidity)
	if err != nil {
		return nil, err
	}
	actualAmountBase, actualAmountQuote, err := pool.CalcActualAmounts(position.LowerTick, position.UpperTick, liquidity)
	if err != nil {
		return nil, err
	}

	amountInDec := math.LegacyNewDecFromInt(amountIn)
	if req.DenomIn == pool.DenomBase {
		if actualAmountBase.IsZero() {
			return nil, types.ErrZeroActualAmountBase
		}
		return &types.QueryCalculationIncreaseLiquidityResponse{
			TokenRequired: sdk.NewCoin(pool.DenomQuote, amountInDec.Mul(actualAmountQuote).Quo(actualAmountBase).TruncateInt()),
		}, nil
	} else if req.DenomIn == pool.DenomQuote {
		if actualAmountQuote.IsZero() {
			return nil, types.ErrZeroActualAmountQuote
		}
		return &types.QueryCalculationIncreaseLiquidityResponse{
			TokenRequired: sdk.NewCoin(pool.DenomBase, amountInDec.Mul(actualAmountBase).Quo(actualAmountQuote).TruncateInt()),
		}, nil
	} else {
		return nil, types.ErrInvalidInDenom
	}
}
