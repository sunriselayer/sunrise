package keeper

import (
	"context"
	"errors"
	"strconv"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors" // aliased by user
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors" // Still needed for PageResponse type in QueryBribeAllocationsResponse
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// BribeAllocations queries BribeAllocations with an optional address filter. Pagination is removed.
func (q queryServer) BribeAllocations(ctx context.Context, req *types.QueryBribeAllocationsRequest) (*types.QueryBribeAllocationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var allocations []types.BribeAllocation
	if req.Address != "" {
		addr, errDecode := q.k.addressCodec.StringToBytes(req.Address)
		if errDecode != nil {
			q.k.Logger().Error("failed to decode address string", "address", req.Address, "error", errDecode)
			return nil, status.Errorf(codes.InvalidArgument, "invalid address format '%s': %v", req.Address, errDecode)
		}

		// var allAllocationsForAddress []types.BribeAllocation // This line seems to be a leftover from pagination logic
		prefix := collections.NewPrefixedTripleRange[sdk.AccAddress, uint64, uint64](addr)
		errWalk := q.k.BribeAllocations.Walk(ctx, prefix, func(key collections.Triple[sdk.AccAddress, uint64, uint64], value types.BribeAllocation) (stop bool, err error) {
			allocations = append(allocations, value)
			return false, nil
		})
		if errWalk != nil {
			q.k.Logger().Error("failed to walk bribe allocations by address", "address", req.Address, "error", errWalk)
			return nil, status.Errorf(codes.Internal, "failed to query bribe allocations by address %s: %v", req.Address, errWalk)
		}
	} else {
		a, err := q.k.GetAllBribeAllocations(ctx)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) && !errors.Is(err, sdkerrors.ErrKeyNotFound) {
				q.k.Logger().Error("failed to get all bribe allocations", "error", err)
				return nil, status.Errorf(codes.Internal, "failed to query all bribe allocations: %v", err)
			}
		}
		if a != nil {
			allocations = a
		}
	}

	// Pagination field in response is nil.
	return &types.QueryBribeAllocationsResponse{BribeAllocations: allocations}, nil
}

// BribeAllocation queries a BribeAllocation by address, epoch ID, and pool ID.
func (q queryServer) BribeAllocation(ctx context.Context, req *types.QueryBribeAllocationRequest) (*types.QueryBribeAllocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request: request is nil")
	}
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request: address is empty")
	}

	addr, err := q.k.addressCodec.StringToBytes(req.Address)
	if err != nil {
		q.k.Logger().Error("failed to decode address string", "address", req.Address, "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid address format '%s': %v", req.Address, err)
	}

	epochId, err := strconv.ParseUint(req.EpochId, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid epoch_id format: %v", err)
	}

	poolId, err := strconv.ParseUint(req.PoolId, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pool_id format: %v", err)
	}

	allocation, err := q.k.GetBribeAllocation(ctx, addr, epochId, poolId)
	if err != nil {
		// Use standard errors.Is for collections.ErrNotFound and sdkerrors.ErrKeyNotFound (if it's a plain error type)
		if errors.Is(err, sdkerrors.ErrKeyNotFound) || errors.Is(err, collections.ErrNotFound) {
			// errorsmod.Wrapf is appropriate here as sdkerrors.ErrKeyNotFound is from cosmos-sdk and we want to wrap it with context.
			return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "bribe allocation not found for address %s, epoch %d, pool %d", req.Address, epochId, poolId)
		}
		q.k.Logger().Error("failed to get bribe allocation", "address", req.Address, "epoch_id", epochId, "pool_id", poolId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get bribe allocation for address %s, epoch %d, pool %d: %v", req.Address, epochId, poolId, err)
	}

	return &types.QueryBribeAllocationResponse{BribeAllocation: allocation}, nil
}
