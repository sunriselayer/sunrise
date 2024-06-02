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
	"github.com/sunriselayer/sunrise/x/swap/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestInFlightPacketQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	msgs := createNInFlightPacket(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetInFlightPacketRequest
		response *types.QueryGetInFlightPacketResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetInFlightPacketRequest{
				SrcPortId:    msgs[0].SrcPortId,
				SrcChannelId: msgs[0].SrcChannelId,
				Sequence:     msgs[0].Sequence,
			},
			response: &types.QueryGetInFlightPacketResponse{InFlightPacket: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetInFlightPacketRequest{
				SrcPortId:    msgs[1].SrcPortId,
				SrcChannelId: msgs[1].SrcChannelId,
				Sequence:     msgs[1].Sequence,
			},
			response: &types.QueryGetInFlightPacketResponse{InFlightPacket: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetInFlightPacketRequest{
				SrcPortId:    strconv.Itoa(100000),
				SrcChannelId: strconv.Itoa(100000),
				Sequence:     uint64(100000),
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
			response, err := keeper.InFlightPacket(ctx, tc.request)
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

func TestInFlightPacketQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	msgs := createNInFlightPacket(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllInFlightPacketRequest {
		return &types.QueryAllInFlightPacketRequest{
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
			resp, err := keeper.InFlightPacketAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.InFlightPacket), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.InFlightPacket),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.InFlightPacketAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.InFlightPacket), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.InFlightPacket),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.InFlightPacketAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.InFlightPacket),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.InFlightPacketAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
