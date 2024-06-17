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

func (k Keeper) CalculationIncreaseLiquidity(ctx context.Context, req *types.QueryCalculationIncreaseLiquidityRequest) (*types.QueryCalculationIncreaseLiquidityResponse, error) {
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

	position, found := k.GetPosition(ctx, req.Id)
	if !found {
		return nil, types.ErrPositionNotFound
	}
	pool, found := k.GetPool(ctx, position.PoolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}
	actualAmountBase, actualAmountQuote, err := pool.CalcActualAmounts(position.LowerTick, position.UpperTick, position.Liquidity)
	if err != nil {
		return nil, err
	}

	amountInDec := math.LegacyNewDecFromInt(amountIn)
	if req.DenomIn == pool.DenomBase {
		return &types.QueryCalculationIncreaseLiquidityResponse{
			TokenRequired: sdk.NewCoin(pool.DenomQuote, amountInDec.Mul(actualAmountBase).Quo(actualAmountQuote).TruncateInt()),
		}, nil
	} else if req.DenomIn == pool.DenomQuote {
		return &types.QueryCalculationIncreaseLiquidityResponse{
			TokenRequired: sdk.NewCoin(pool.DenomBase, amountInDec.Mul(actualAmountQuote).Quo(actualAmountBase).TruncateInt()),
		}, nil
	} else {
		return nil, types.ErrInvalidInDenom
	}
}
