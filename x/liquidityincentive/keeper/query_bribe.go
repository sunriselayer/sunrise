package keeper

import (
	"context"
	"errors"
	"strconv"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// Bribes queries Bribes with optional filters. Pagination is removed.
// Filters for epoch_id and pool_id are applied if provided in the request as non-empty strings.
func (q queryServer) Bribes(ctx context.Context, req *types.QueryBribesRequest) (*types.QueryBribesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var bribes []types.Bribe
	var err error

	applyEpochIdFilter := req.EpochId != ""
	applyPoolIdFilter := req.PoolId != ""

	var epochId uint64
	if applyEpochIdFilter {
		epochId, err = strconv.ParseUint(req.EpochId, 10, 64)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid epoch_id format '%s': %v", req.EpochId, err)
		}
	}

	var poolId uint64
	if applyPoolIdFilter {
		poolId, err = strconv.ParseUint(req.PoolId, 10, 64)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid pool_id format '%s': %v", req.PoolId, err)
		}
	}

	if applyEpochIdFilter && applyPoolIdFilter {
		b, errGet := q.k.GetBribeByEpochAndPool(ctx, epochId, poolId)
		if errGet != nil {
			if !errors.Is(errGet, collections.ErrNotFound) && !errors.Is(errGet, sdkerrors.ErrKeyNotFound) {
				q.k.Logger().Error("failed to get bribes by epoch and pool id", "epoch_id", epochId, "pool_id", poolId, "error", errGet)
				return nil, status.Errorf(codes.Internal, "failed to query bribes by epoch id %d and pool id %d: %v", epochId, poolId, errGet)
			}
		}
		if b != nil {
			bribes = b
		}
	} else if applyEpochIdFilter {
		bribes, err = q.k.GetAllBribeByEpochId(ctx, epochId)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) && !errors.Is(err, sdkerrors.ErrKeyNotFound) {
				q.k.Logger().Error("failed to get bribes by epoch id", "epoch_id", epochId, "error", err)
				return nil, status.Errorf(codes.Internal, "failed to query bribes by epoch id %d: %v", epochId, err)
			}
		}
	} else if applyPoolIdFilter {
		bribes, err = q.k.GetAllBribeByPoolId(ctx, poolId)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) && !errors.Is(err, sdkerrors.ErrKeyNotFound) {
				q.k.Logger().Error("failed to get bribes by pool id", "pool_id", poolId, "error", err)
				return nil, status.Errorf(codes.Internal, "failed to query bribes by pool id %d: %v", poolId, err)
			}
		}
	} else {
		bribes, err = q.k.GetAllBribes(ctx)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) && !errors.Is(err, sdkerrors.ErrKeyNotFound) {
				q.k.Logger().Error("failed to get all bribes", "error", err)
				return nil, status.Errorf(codes.Internal, "failed to query all bribes: %v", err)
			}
		}
	}

	return &types.QueryBribesResponse{Bribes: bribes}, nil
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
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "bribe with id %d not found", req.Id)
	}

	return &types.QueryBribeResponse{Bribe: bribe}, nil
}
