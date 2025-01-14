package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (q queryServer) Epochs(ctx context.Context, req *types.QueryEpochsRequest) (*types.QueryEpochsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	epochs, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Epochs,
		req.Pagination,
		func(key uint64, value types.Epoch) (types.Epoch, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryEpochsResponse{Epochs: epochs, Pagination: pageRes}, nil
}

func (q queryServer) Epoch(ctx context.Context, req *types.QueryEpochRequest) (*types.QueryEpochResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	epoch, found := q.k.GetEpoch(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryEpochResponse{Epoch: epoch}, nil
}
