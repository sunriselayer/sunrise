package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestGaugeQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LiquidityincentiveKeeper(t)
	msgs := createNGauge(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetGaugeRequest
		response *types.QueryGetGaugeResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetGaugeRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetGaugeResponse{Gauge: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetGaugeRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetGaugeResponse{Gauge: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetGaugeRequest{
				Index: strconv.Itoa(100000),
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
			response, err := keeper.Gauge(ctx, tc.request)
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
	keeper, ctx := keepertest.LiquidityincentiveKeeper(t)
	msgs := createNGauge(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllGaugeRequest {
		return &types.QueryAllGaugeRequest{
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
			resp, err := keeper.GaugeAll(ctx, request(nil, uint64(i), uint64(step), false))
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
			resp, err := keeper.GaugeAll(ctx, request(next, 0, uint64(step), false))
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
		resp, err := keeper.GaugeAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Gauge),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.GaugeAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
