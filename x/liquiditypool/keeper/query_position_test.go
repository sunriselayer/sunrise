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
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestPositionQuerySingle(t *testing.T) {
	keeper, _, ctx := keepertest.LiquiditypoolKeeper(t)
	msgs := createNPosition(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryPositionRequest
		response *types.QueryPositionResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryPositionRequest{Id: msgs[0].Id},
			response: &types.QueryPositionResponse{Position: types.PositionInfo{Position: msgs[0]}},
		},
		{
			desc:     "Second",
			request:  &types.QueryPositionRequest{Id: msgs[1].Id},
			response: &types.QueryPositionResponse{Position: types.PositionInfo{Position: msgs[1]}},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryPositionRequest{Id: uint64(len(msgs) + 1)},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Position(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response.Position),
					nullify.Fill(response.Position),
				)
			}
		})
	}
}

func TestPositionQueryPaginated(t *testing.T) {
	keeper, _, ctx := keepertest.LiquiditypoolKeeper(t)
	msgs := createNPosition(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryPositionsRequest {
		return &types.QueryPositionsRequest{
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
			resp, err := keeper.Positions(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Positions), step)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.Positions(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Positions), step)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.Positions(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.Positions(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
