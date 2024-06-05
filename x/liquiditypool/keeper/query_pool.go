package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

<<<<<<< HEAD
func (k Keeper) WrapPoolInfo(ctx context.Context, pool types.Pool) types.PoolInfo {
	return types.PoolInfo{
		Pool:        pool,
		AmountBase:  sdk.Coin{}, // TODO: need to add liquidity tracker in pool
		AmountQuote: sdk.Coin{}, // TODO: need to add liquidity tracker in pool
	}
}

func (k Keeper) PoolAll(ctx context.Context, req *types.QueryAllPoolRequest) (*types.QueryAllPoolResponse, error) {
=======
func (k Keeper) Pools(ctx context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
>>>>>>> 3988f81665ee85f01cc5adb24fcb7beb3e5cb010
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pools []types.PoolInfo

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	poolStore := prefix.NewStore(store, types.KeyPrefix(types.PoolKey))

	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
		var pool types.Pool
		if err := k.cdc.Unmarshal(value, &pool); err != nil {
			return err
		}

		pools = append(pools, k.WrapPoolInfo(ctx, pool))
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolsResponse{Pools: pools, Pagination: pageRes}, nil
}

func (k Keeper) Pool(ctx context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	pool, found := k.GetPool(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

<<<<<<< HEAD
	return &types.QueryGetPoolResponse{Pool: k.WrapPoolInfo(ctx, pool)}, nil
=======
	return &types.QueryPoolResponse{Pool: pool}, nil
>>>>>>> 3988f81665ee85f01cc5adb24fcb7beb3e5cb010
}
