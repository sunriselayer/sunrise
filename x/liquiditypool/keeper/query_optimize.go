package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OptimizeSwapExactAmountIn(ctx context.Context, req *types.QueryOptimizeSwapExactAmountInRequest) (*types.QueryOptimizeSwapExactAmountInResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var err error

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOptimizeSwapExactAmountInResponse{}, nil
}

func (k Keeper) OptimizeSwapExactAmountOut(ctx context.Context, req *types.QueryOptimizeSwapExactAmountOutRequest) (*types.QueryOptimizeSwapExactAmountOutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var err error

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOptimizeSwapExactAmountOutResponse{}, nil
}
