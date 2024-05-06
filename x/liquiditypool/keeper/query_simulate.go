package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SimulateSwapExactAmountIn(ctx context.Context, req *types.QuerySimulateSwapExactAmountInRequest) (*types.QuerySimulateSwapExactAmountInResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	tokensVia, tokenOut, err := k.SwapExactAmountInMultiRoute(ctx, req.Routes, req.TokenIn, true, nil, nil)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySimulateSwapExactAmountInResponse{
		TokensVia: tokensVia,
		TokenOut:  *tokenOut,
	}, nil
}

func (k Keeper) SimulateSwapExactAmountOut(ctx context.Context, req *types.QuerySimulateSwapExactAmountOutRequest) (*types.QuerySimulateSwapExactAmountOutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	tokenIn, err := k.SwapExactAmountOutMultiPool(ctx, req.Route, req.TokenOut, true, nil, nil)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySimulateSwapExactAmountOutResponse{
		TokenIn: *tokenIn,
	}, nil
}

func (k Keeper) SimulateJoinPool(ctx context.Context, req *types.QuerySimulateJoinPoolRequest) (*types.QuerySimulateJoinPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	shareAmount, err := k.joinPool(ctx, req.PoolId, req.BaseToken, req.QuoteToken, true, nil, nil)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySimulateJoinPoolResponse{
		ShareAmount: *shareAmount,
	}, nil
}

func (k Keeper) SimulateExitPool(ctx context.Context, req *types.QuerySimulateExitPoolRequest) (*types.QuerySimulateExitPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	tokensOut, err := k.exitPool(ctx, req.PoolId, req.ShareAmount, true, nil, nil, nil)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySimulateExitPoolResponse{
		TokensOut: tokensOut,
	}, nil
}
