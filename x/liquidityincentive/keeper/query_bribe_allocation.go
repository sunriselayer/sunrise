package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// BribeAllocations queries all BribeAllocations with pagination.
func (q queryServer) BribeAllocations(ctx context.Context, req *types.QueryBribeAllocationsRequest) (*types.QueryBribeAllocationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	allocations, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.BribeAllocations,
		req.Pagination,
		func(key collections.Triple[sdk.AccAddress, uint64, uint64], value types.BribeAllocation) (types.BribeAllocation, error) {
			return value, nil
		},
	)
	if err != nil {
		q.k.Logger().Error("failed to paginate bribe allocations", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to query bribe allocations: %v", err)
	}

	return &types.QueryBribeAllocationsResponse{BribeAllocations: allocations, Pagination: pageRes}, nil
}

// BribeAllocationsByAddress queries BribeAllocations associated with a specific address.
// Note: Pagination is not implemented as it's not defined in the proto request.
func (q queryServer) BribeAllocationsByAddress(ctx context.Context, req *types.QueryBribeAllocationsByAddressRequest) (*types.QueryBribeAllocationsByAddressResponse, error) {
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

	var allocations []types.BribeAllocation
	prefix := collections.NewPrefixedTripleRange[sdk.AccAddress, uint64, uint64](addr)
	// Walk through relevant allocations
	err = q.k.BribeAllocations.Walk(ctx, prefix, func(key collections.Triple[sdk.AccAddress, uint64, uint64], value types.BribeAllocation) (stop bool, err error) {
		allocations = append(allocations, value)
		return false, nil
	})
	if err != nil {
		q.k.Logger().Error("failed to walk bribe allocations by address", "address", req.Address, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to query bribe allocations by address %s: %v", req.Address, err)
	}

	return &types.QueryBribeAllocationsByAddressResponse{BribeAllocations: allocations}, nil
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

	allocation, err := q.k.GetBribeAllocation(ctx, addr, req.EpochId, req.PoolId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errors.Wrapf(sdkerrors.ErrKeyNotFound, "bribe allocation not found for address %s, epoch %d, pool %d", req.Address, req.EpochId, req.PoolId)
		}
		q.k.Logger().Error("failed to get bribe allocation", "address", req.Address, "epoch_id", req.EpochId, "pool_id", req.PoolId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get bribe allocation for address %s, epoch %d, pool %d: %v", req.Address, req.EpochId, req.PoolId, err)
	}

	return &types.QueryBribeAllocationResponse{BribeAllocation: allocation}, nil
}
