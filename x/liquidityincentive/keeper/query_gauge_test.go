package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestGaugeQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNGauge(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGaugeRequest
		response *types.QueryGaugeResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGaugeRequest{
				PreviousEpochId: msgs[0].PreviousEpochId,
				PoolId:          msgs[0].PoolId,
			},
			response: &types.QueryGaugeResponse{Gauge: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGaugeRequest{
				PreviousEpochId: msgs[1].PreviousEpochId,
				PoolId:          msgs[1].PoolId,
			},
			response: &types.QueryGaugeResponse{Gauge: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGaugeRequest{
				PreviousEpochId: 100000,
				PoolId:          100000,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.Gauge(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestGaugeQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNGauge(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryGaugesRequest {
		return &types.QueryGaugesRequest{
			PreviousEpochId: 1,
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.Gauges(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Gauge), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Gauge),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.Gauges(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Gauge), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Gauge),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.Gauges(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Gauge),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.Gauges(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
