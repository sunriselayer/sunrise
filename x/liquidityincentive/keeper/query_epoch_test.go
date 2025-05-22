package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	t.Run("EmptyResults", func(t *testing.T) {
		// Query with offset beyond available epochs
		resp, err := qs.Epochs(f.ctx, request(nil, uint64(len(msgs)+1), 10, false))
		require.NoError(t, err)
		require.Empty(t, resp.Epochs)
	})
	t.Run("InvalidPagination", func(t *testing.T) {
		// Test with invalid pagination parameters
		req := &types.QueryEpochsRequest{
			Pagination: &query.PageRequest{
				Offset: 999999, // Very large offset
				Limit:  0,      // Invalid limit
			},
		}
		resp, err := qs.Epochs(f.ctx, req)
		require.NoError(t, err)
		require.Empty(t, resp.Epochs)
	})
}

func TestEpochQueryByBlockHeight(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)

	// Create epochs with specific block height ranges
	epochs := []types.Epoch{
		{
			Id:         1,
			StartBlock: 100,
			EndBlock:   200,
			Gauges:     []types.Gauge{},
		},
		{
			Id:         2,
			StartBlock: 201,
			EndBlock:   300,
			Gauges:     []types.Gauge{},
		},
	}

	for _, epoch := range epochs {
		err := f.keeper.SetEpoch(f.ctx, epoch)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		blockHeight   int64
		expectedEpoch *types.Epoch
		expectError   bool
	}{
		{
			name:          "BlockInFirstEpoch",
			blockHeight:   150,
			expectedEpoch: &epochs[0],
		},
		{
			name:          "BlockInSecondEpoch",
			blockHeight:   250,
			expectedEpoch: &epochs[1],
		},
		{
			name:        "BlockBeforeAnyEpoch",
			blockHeight: 50,
			expectError: true,
		},
		{
			name:        "BlockAfterAnyEpoch",
			blockHeight: 350,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set the block height in the context
			sdkCtx := sdk.UnwrapSDKContext(f.ctx)
			ctx := sdkCtx.WithBlockHeight(tc.blockHeight)

			// Get all epochs
			resp, err := qs.Epochs(ctx, &types.QueryEpochsRequest{})
			require.NoError(t, err)

			if tc.expectError {
				// For blocks outside epoch ranges, we should find no matching epochs
				found := false
				for _, epoch := range resp.Epochs {
					if epoch.StartBlock <= tc.blockHeight && tc.blockHeight <= epoch.EndBlock {
						found = true
						break
					}
				}
				require.False(t, found, "Found epoch for block height %d when none should exist", tc.blockHeight)
			} else {
				require.NotEmpty(t, resp.Epochs)
				found := false
				for _, epoch := range resp.Epochs {
					if epoch.StartBlock <= tc.blockHeight && tc.blockHeight <= epoch.EndBlock {
						found = true
						require.Equal(t, tc.expectedEpoch.Id, epoch.Id)
						break
					}
				}
				require.True(t, found, "No epoch found for block height %d", tc.blockHeight)
			}
		})
	}
}
