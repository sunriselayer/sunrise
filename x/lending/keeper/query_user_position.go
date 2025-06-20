package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func (q queryServer) UserPosition(ctx context.Context, req *types.QueryUserPositionRequest) (*types.QueryUserPositionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.UserAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "user address cannot be empty")
	}

	if req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "denom cannot be empty")
	}

	position, err := q.k.UserPositions.Get(ctx, collections.Join(req.UserAddress, req.Denom))
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "position not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryUserPositionResponse{Position: position}, nil
}

func (q queryServer) UserPositions(ctx context.Context, req *types.QueryUserPositionsRequest) (*types.QueryUserPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.UserAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "user address cannot be empty")
	}

	positions, pageResp, err := query.CollectionFilteredPaginate(
		ctx,
		q.k.UserPositions,
		req.Pagination,
		func(key collections.Pair[string, string], value types.UserPosition) (include bool, err error) {
			// Only include positions for the requested user
			return key.K1() == req.UserAddress, nil
		},
		func(key collections.Pair[string, string], value types.UserPosition) (types.UserPosition, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryUserPositionsResponse{
		Positions:  positions,
		Pagination: pageResp,
	}, nil
}