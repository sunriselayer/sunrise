package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func TestEpochQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.LiquidityincentiveKeeper(t)
	msgs := createNEpoch(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetEpochRequest
		response *types.QueryGetEpochResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetEpochRequest{Id: msgs[0].Id},
			response: &types.QueryGetEpochResponse{Epoch: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetEpochRequest{Id: msgs[1].Id},
			response: &types.QueryGetEpochResponse{Epoch: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetEpochRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Epoch(ctx, tc.request)
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

func TestEpochQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.LiquidityincentiveKeeper(t)
	msgs := createNEpoch(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllEpochRequest {
		return &types.QueryAllEpochRequest{
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
			resp, err := keeper.EpochAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Epoch), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Epoch),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.EpochAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Epoch), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Epoch),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.EpochAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Epoch),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.EpochAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
