package keeper_test

import (
    "strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
	"github.com/sunriselayer/sunrise-app/testutil/nullify"
	keepertest "github.com/sunriselayer/sunrise-app/testutil/keeper"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestPairQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	msgs := createNPair(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetPairRequest
		response *types.QueryGetPairResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetPairRequest{
			    Index: msgs[0].Index,
                
			},
			response: &types.QueryGetPairResponse{Pair: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetPairRequest{
			    Index: msgs[1].Index,
                
			},
			response: &types.QueryGetPairResponse{Pair: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetPairRequest{
			    Index:strconv.Itoa(100000),
                
			},
			err:     status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Pair(ctx, tc.request)
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

func TestPairQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	msgs := createNPair(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPairRequest {
		return &types.QueryAllPairRequest{
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
			resp, err := keeper.PairAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Pair), step)
			require.Subset(t,
            	nullify.Fill(msgs),
            	nullify.Fill(resp.Pair),
            )
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.PairAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Pair), step)
			require.Subset(t,
            	nullify.Fill(msgs),
            	nullify.Fill(resp.Pair),
            )
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.PairAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Pair),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.PairAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
