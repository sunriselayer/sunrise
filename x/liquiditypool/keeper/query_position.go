package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Positions(ctx context.Context, req *types.QueryPositionsRequest) (*types.QueryPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var positions []types.Position

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	positionStore := prefix.NewStore(store, types.KeyPrefix(types.PositionKey))

	pageRes, err := query.Paginate(positionStore, req.Pagination, func(key []byte, value []byte) error {
		var position types.Position
		if err := k.cdc.Unmarshal(value, &position); err != nil {
			return err
		}

		positions = append(positions, position)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPositionsResponse{Positions: positions, Pagination: pageRes}, nil
}

func (k Keeper) Position(ctx context.Context, req *types.QueryPositionRequest) (*types.QueryPositionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	position, found := k.GetPosition(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryPositionResponse{Position: position}, nil
}

func (k Keeper) AddressPositions(ctx context.Context, req *types.QueryAddressPositionsRequest) (*types.QueryAddressPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var positions []types.Position

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	positionStore := prefix.NewStore(store, types.KeyPrefix(types.PositionKey))

	pageRes, err := query.Paginate(positionStore, req.Pagination, func(key []byte, value []byte) error {
		var position types.Position
		if err := k.cdc.Unmarshal(value, &position); err != nil {
			return err
		}

		if position.Address == req.Address {
			positions = append(positions, position)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAddressPositionsResponse{Positions: positions, Pagination: pageRes}, nil
}

func (k Keeper) PositionFees(ctx context.Context, req *types.QueryPositionFeesRequest) (*types.QueryPositionFeesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// todo: implement logic

	return &types.QueryPositionFeesResponse{}, nil
}
