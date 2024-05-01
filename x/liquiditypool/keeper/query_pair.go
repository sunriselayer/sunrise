package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PairAll(ctx context.Context, req *types.QueryAllPairRequest) (*types.QueryAllPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pairs []types.Pair

    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	pairStore := prefix.NewStore(store, types.KeyPrefix(types.PairKeyPrefix))

	pageRes, err := query.Paginate(pairStore, req.Pagination, func(key []byte, value []byte) error {
		var pair types.Pair
		if err := k.cdc.Unmarshal(value, &pair); err != nil {
			return err
		}

		pairs = append(pairs, pair)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPairResponse{Pair: pairs, Pagination: pageRes}, nil
}

func (k Keeper) Pair(ctx context.Context, req *types.QueryGetPairRequest) (*types.QueryGetPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetPair(
	    ctx,
	    req.Index,
        )
	if !found {
	    return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPairResponse{Pair: val}, nil
}