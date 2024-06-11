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

func (k Keeper) GaugeAll(ctx context.Context, req *types.QueryAllGaugeRequest) (*types.QueryAllGaugeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var gauges []types.Gauge

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	gaugeStore := prefix.NewStore(store, types.KeyPrefix(types.GaugeKeyPrefix))

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

	return &types.QueryAllGaugeResponse{Gauge: gauges, Pagination: pageRes}, nil
}

func (k Keeper) Gauge(ctx context.Context, req *types.QueryGetGaugeRequest) (*types.QueryGetGaugeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetGauge(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetGaugeResponse{Gauge: val}, nil
}
