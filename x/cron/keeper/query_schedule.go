package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/cron/types"
)

func (q queryServer) Schedules(c context.Context, req *types.QuerySchedulesRequest) (*types.QuerySchedulesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	schedules, pageRes, err := query.CollectionPaginate(ctx, q.k.Schedules, req.Pagination, func(_ string, value types.Schedule) (types.Schedule, error) {
		return value, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySchedulesResponse{Schedules: schedules, Pagination: pageRes}, nil
}

func (q queryServer) Schedule(c context.Context, req *types.QueryGetScheduleRequest) (*types.QueryGetScheduleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := q.k.GetSchedule(
		ctx,
		req.Name,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "schedule not found")
	}

	return &types.QueryGetScheduleResponse{Schedule: *val}, nil
}
