package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PositionsIncentives(goCtx context.Context, req *types.QueryPositionsIncentivesRequest) (*types.QueryPositionsIncentivesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	// TODO: implement

	return &types.QueryPositionsIncentivesResponse{}, nil
}

func (k Keeper) PositionIncentives(goCtx context.Context, req *types.QueryPositionIncentivesRequest) (*types.QueryPositionIncentivesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	// TODO: implement

	return &types.QueryPositionIncentivesResponse{}, nil
}

func (k Keeper) AddressIncentives(goCtx context.Context, req *types.QueryAddressIncentivesRequest) (*types.QueryAddressIncentivesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	// TODO: implement

	return &types.QueryAddressIncentivesResponse{}, nil
}
