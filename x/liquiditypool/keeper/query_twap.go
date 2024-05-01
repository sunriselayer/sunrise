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

func (k Keeper) TwapAll(ctx context.Context, req *types.QueryAllTwapRequest) (*types.QueryAllTwapResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var twaps []types.Twap

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	twapStore := prefix.NewStore(store, types.KeyPrefix(types.TwapKeyPrefix))

	pageRes, err := query.Paginate(twapStore, req.Pagination, func(key []byte, value []byte) error {
		var twap types.Twap
		if err := k.cdc.Unmarshal(value, &twap); err != nil {
			return err
		}

		twaps = append(twaps, twap)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTwapResponse{Twap: twaps, Pagination: pageRes}, nil
}

func (k Keeper) Twap(ctx context.Context, req *types.QueryGetTwapRequest) (*types.QueryGetTwapResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetTwap(
		ctx,
		req.BaseDenom,
		req.QuoteDenom,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetTwapResponse{Twap: val}, nil
}
