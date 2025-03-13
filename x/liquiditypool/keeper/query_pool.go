package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k Keeper) WrapPoolInfo(ctx context.Context, pool types.Pool) types.PoolInfo {
	tokenBase := k.bankKeeper.GetBalance(ctx, pool.GetAddress(), pool.DenomBase)
	tokenQuote := k.bankKeeper.GetBalance(ctx, pool.GetAddress(), pool.DenomQuote)
	return types.PoolInfo{
		Pool:       pool,
		TokenBase:  tokenBase,
		TokenQuote: tokenQuote,
	}
}

func (q queryServer) Pools(ctx context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	pools, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Pools,
		req.Pagination,
		func(_ uint64, value types.Pool) (types.PoolInfo, error) {
			return q.k.WrapPoolInfo(ctx, value), nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolsResponse{Pools: pools, Pagination: pageRes}, nil
}

func (q queryServer) Pool(ctx context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	pool, found, err := q.k.GetPool(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryPoolResponse{Pool: q.k.WrapPoolInfo(ctx, pool)}, nil
}
