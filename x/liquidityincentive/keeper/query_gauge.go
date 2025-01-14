package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (q queryServer) Gauges(ctx context.Context, req *types.QueryGaugesRequest) (*types.QueryGaugesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	gauges, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Gauges,
		req.Pagination,
		func(key collections.Pair[uint64, uint64], value types.Gauge) (types.Gauge, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGaugesResponse{Gauge: gauges, Pagination: pageRes}, nil
}

func (q queryServer) Gauge(ctx context.Context, req *types.QueryGaugeRequest) (*types.QueryGaugeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := q.k.GetGauge(
		ctx,
		req.PreviousEpochId,
		req.PoolId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGaugeResponse{Gauge: val}, nil
}
