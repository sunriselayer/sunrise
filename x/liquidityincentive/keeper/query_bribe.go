package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// Bribes queries all Bribes with pagination.
func (q queryServer) Bribes(ctx context.Context, req *types.QueryBribesRequest) (*types.QueryBribesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	bribes, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Bribes,
		req.Pagination,
		func(key uint64, value types.Bribe) (types.Bribe, error) {
			return value, nil
		},
	)
	if err != nil {
		q.k.Logger().Error("failed to paginate bribes", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to query bribes: %v", err)
	}

	return &types.QueryBribesResponse{Bribes: bribes, Pagination: pageRes}, nil
}

// Bribe queries a Bribe by its ID.
func (q queryServer) Bribe(ctx context.Context, req *types.QueryBribeRequest) (*types.QueryBribeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	bribe, found, err := q.k.GetBribe(ctx, req.Id)
	if err != nil {
		q.k.Logger().Error("failed to get bribe", "id", req.Id, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get bribe with id %d: %v", req.Id, err)
	}
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrKeyNotFound, "bribe with id %d not found", req.Id)
	}

	return &types.QueryBribeResponse{Bribe: bribe}, nil
}

// BribesByEpochId queries all Bribes associated with a specific epoch ID.
// Note: Pagination is not supported as it's not defined in the proto request.
func (q queryServer) BribesByEpochId(ctx context.Context, req *types.QueryBribesByEpochIdRequest) (*types.QueryBribesByEpochIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Retrieve all bribes for the given epoch ID from the keeper (no pagination)
	bribes, err := q.k.GetAllBribeByEpochId(ctx, req.EpochId)
	if err != nil {
		q.k.Logger().Error("failed to get bribes by epoch id", "epoch_id", req.EpochId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to query bribes by epoch id %d: %v", req.EpochId, err)
	}

	// Create response (no pagination)
	return &types.QueryBribesByEpochIdResponse{Bribes: bribes}, nil
}

// BribesByPoolId queries all Bribes associated with a specific pool ID.
// Note: Pagination is not supported as it's not defined in the proto request.
func (q queryServer) BribesByPoolId(ctx context.Context, req *types.QueryBribesByPoolIdRequest) (*types.QueryBribesByPoolIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Retrieve all bribes for the given pool ID from the keeper (no pagination)
	bribes, err := q.k.GetAllBribeByPoolId(ctx, req.PoolId)
	if err != nil {
		q.k.Logger().Error("failed to get bribes by pool id", "pool_id", req.PoolId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to query bribes by pool id %d: %v", req.PoolId, err)
	}

	// Create response (no pagination)
	return &types.QueryBribesByPoolIdResponse{Bribes: bribes}, nil
}

// BribesByEpochAndPoolId queries all Bribes associated with a specific epoch ID and pool ID.
// Note: Pagination is not supported as it's not defined in the proto request.
func (q queryServer) BribesByEpochAndPoolId(ctx context.Context, req *types.QueryBribesByEpochAndPoolIdRequest) (*types.QueryBribesByEpochAndPoolIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Retrieve all bribes for the given epoch and pool ID from the keeper (no pagination)
	bribes, err := q.k.GetBribeByEpochAndPool(ctx, req.EpochId, req.PoolId)
	if err != nil {
		q.k.Logger().Error("failed to get bribes by epoch and pool id", "epoch_id", req.EpochId, "pool_id", req.PoolId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to query bribes by epoch id %d and pool id %d: %v", req.EpochId, req.PoolId, err)
	}

	// Create response (no pagination)
	return &types.QueryBribesByEpochAndPoolIdResponse{Bribes: bribes}, nil
}
