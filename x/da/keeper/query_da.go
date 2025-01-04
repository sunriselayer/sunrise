package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) PublishedData(goCtx context.Context, req *types.QueryPublishedDataRequest) (*types.QueryPublishedDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	data := q.k.GetPublishedData(ctx, req.MetadataUri)
	return &types.QueryPublishedDataResponse{Data: data}, nil
}

func (q queryServer) AllPublishedData(goCtx context.Context, req *types.QueryAllPublishedDataRequest) (*types.QueryAllPublishedDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAllPublishedDataResponse{Data: q.k.GetAllPublishedData(ctx)}, nil
}

func (q queryServer) ZkpProofThreshold(goCtx context.Context, req *types.QueryZkpProofThresholdRequest) (*types.QueryZkpProofThresholdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryZkpProofThresholdResponse{Threshold: q.k.GetZkpThreshold(ctx, req.ShardCount)}, nil
}
