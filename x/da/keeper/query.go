package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) PublishedData(goCtx context.Context, req *types.QueryPublishedDataRequest) (*types.QueryPublishedDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	data := k.GetPublishedData(ctx, req.MetadataUri)
	return &types.QueryPublishedDataResponse{Data: data}, nil
}

func (k Keeper) AllPublishedData(goCtx context.Context, req *types.QueryAllPublishedDataRequest) (*types.QueryAllPublishedDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAllPublishedDataResponse{Data: k.GetAllPublishedData(ctx)}, nil
}

func (k Keeper) ZkpProofThreshold(goCtx context.Context, req *types.QueryZkpProofThresholdRequest) (*types.QueryZkpProofThresholdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryZkpProofThresholdResponse{Threshold: k.GetZkpThreshold(ctx, req.ShardCount)}, nil
}
