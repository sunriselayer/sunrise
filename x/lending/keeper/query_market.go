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

func (q queryServer) Market(ctx context.Context, req *types.QueryMarketRequest) (*types.QueryMarketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "denom cannot be empty")
	}

	market, err := q.k.Markets.Get(ctx, req.Denom)
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "market not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMarketResponse{Market: market}, nil
}

func (q queryServer) Markets(ctx context.Context, req *types.QueryMarketsRequest) (*types.QueryMarketsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	markets, pageResp, err := query.CollectionPaginate(
		ctx,
		q.k.Markets,
		req.Pagination,
		func(key string, value types.Market) (types.Market, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMarketsResponse{
		Markets:    markets,
		Pagination: pageResp,
	}, nil
}