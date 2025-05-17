package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func TestEpochQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNEpoch(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryEpochRequest
		response *types.QueryEpochResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryEpochRequest{Id: msgs[0].Id},
			response: &types.QueryEpochResponse{Epoch: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryEpochRequest{Id: msgs[1].Id},
			response: &types.QueryEpochResponse{Epoch: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryEpochRequest{Id: 9999},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.Epoch(f.ctx, tc.request)
			if tc.err != nil {
				if tc.desc == "KeyNotFound" {
					require.EqualError(t, err, tc.err.Error())
				} else {
					require.ErrorIs(t, err, tc.err)
				}
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
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNEpoch(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryEpochsRequest {
		return &types.QueryEpochsRequest{
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
			resp, err := qs.Epochs(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Epochs), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Epochs),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.Epochs(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Epochs), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Epochs),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.Epochs(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Epochs),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.Epochs(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
