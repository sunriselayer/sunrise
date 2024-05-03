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
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestTwapQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	msgs := createNTwap(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetTwapRequest
		response *types.QueryGetTwapResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetTwapRequest{
				BaseDenom:  msgs[0].BaseDenom,
				QuoteDenom: msgs[0].QuoteDenom,
			},
			response: &types.QueryGetTwapResponse{Twap: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetTwapRequest{
				BaseDenom:  msgs[1].BaseDenom,
				QuoteDenom: msgs[1].QuoteDenom,
			},
			response: &types.QueryGetTwapResponse{Twap: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetTwapRequest{
				BaseDenom:  "base" + strconv.Itoa(100000),
				QuoteDenom: "quote" + strconv.Itoa(100000),
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
			response, err := keeper.Twap(ctx, tc.request)
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

func TestTwapQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	msgs := createNTwap(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllTwapRequest {
		return &types.QueryAllTwapRequest{
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
			resp, err := keeper.TwapAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Twap), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Twap),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.TwapAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Twap), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Twap),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.TwapAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Twap),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.TwapAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
