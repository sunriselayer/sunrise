package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Gauges(ctx context.Context, req *types.QueryGaugesRequest) (*types.QueryGaugesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var gauges []types.Gauge

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	gaugeStore := prefix.NewStore(store, types.GaugeKeyPrefixByPreviousEpochId(req.PreviousEpochId))

	pageRes, err := query.Paginate(gaugeStore, req.Pagination, func(key []byte, value []byte) error {
		var gauge types.Gauge
		if err := k.cdc.Unmarshal(value, &gauge); err != nil {
			return err
		}

		gauges = append(gauges, gauge)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGaugesResponse{Gauge: gauges, Pagination: pageRes}, nil
}

func (k Keeper) Gauge(ctx context.Context, req *types.QueryGaugeRequest) (*types.QueryGaugeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetGauge(
		ctx,
		req.PreviousEpochId,
		req.PoolId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGaugeResponse{Gauge: val}, nil
}
