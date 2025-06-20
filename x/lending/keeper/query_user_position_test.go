package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/lending/keeper"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func TestUserPositionQuery(t *testing.T) {
	f := initFixture(t)
	defer f.ctrl.Finish()

	queryServer := keeper.NewQueryServerImpl(f.keeper)

	// Create test addresses
	user1 := sdk.AccAddress("user1").String()
	user2 := sdk.AccAddress("user2").String()

	// Create test positions
	positions := []types.UserPosition{
		{
			UserAddress:      user1,
			Denom:            "usdc",
			Amount:           math.NewInt(100000),
			LastRewardIndex:  math.LegacyOneDec(),
		},
		{
			UserAddress:      user1,
			Denom:            "atom",
			Amount:           math.NewInt(50000),
			LastRewardIndex:  math.LegacyNewDecWithPrec(95, 2), // 0.95
		},
		{
			UserAddress:      user2,
			Denom:            "usdc",
			Amount:           math.NewInt(200000),
			LastRewardIndex:  math.LegacyNewDecWithPrec(11, 1), // 1.1
		},
	}

	for _, pos := range positions {
		err := f.keeper.SetUserPosition(f.ctx, pos)
		require.NoError(t, err)
	}

	// Test Query single user position
	t.Run("query single user position", func(t *testing.T) {
		resp, err := queryServer.UserPosition(f.ctx, &types.QueryUserPositionRequest{
			UserAddress: user1,
			Denom:       "usdc",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, user1, resp.Position.UserAddress)
		require.Equal(t, "usdc", resp.Position.Denom)
		require.Equal(t, math.NewInt(100000), resp.Position.Amount)
		require.Equal(t, math.LegacyOneDec(), resp.Position.LastRewardIndex)
	})

	// Test Query non-existent position
	t.Run("query non-existent position", func(t *testing.T) {
		_, err := queryServer.UserPosition(f.ctx, &types.QueryUserPositionRequest{
			UserAddress: user1,
			Denom:       "nonexistent",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "position not found")
	})

	// Test Query all positions for a user
	t.Run("query all positions for user1", func(t *testing.T) {
		resp, err := queryServer.UserPositions(f.ctx, &types.QueryUserPositionsRequest{
			UserAddress: user1,
			Pagination: &query.PageRequest{
				Limit: 10,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Positions, 2)
		
		// Check that we only got user1's positions
		for _, pos := range resp.Positions {
			require.Equal(t, user1, pos.UserAddress)
		}
		
		// Check denoms
		denoms := []string{resp.Positions[0].Denom, resp.Positions[1].Denom}
		require.Contains(t, denoms, "usdc")
		require.Contains(t, denoms, "atom")
	})

	// Test Query all positions for user2
	t.Run("query all positions for user2", func(t *testing.T) {
		resp, err := queryServer.UserPositions(f.ctx, &types.QueryUserPositionsRequest{
			UserAddress: user2,
			Pagination: &query.PageRequest{
				Limit: 10,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Positions, 1)
		require.Equal(t, user2, resp.Positions[0].UserAddress)
		require.Equal(t, "usdc", resp.Positions[0].Denom)
	})

	// Test pagination
	t.Run("query positions with pagination", func(t *testing.T) {
		// Query first page with limit 1
		resp, err := queryServer.UserPositions(f.ctx, &types.QueryUserPositionsRequest{
			UserAddress: user1,
			Pagination: &query.PageRequest{
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Positions, 1)
		require.NotNil(t, resp.Pagination)
		require.NotEmpty(t, resp.Pagination.NextKey)

		// Query second page
		resp2, err := queryServer.UserPositions(f.ctx, &types.QueryUserPositionsRequest{
			UserAddress: user1,
			Pagination: &query.PageRequest{
				Key:   resp.Pagination.NextKey,
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp2)
		require.Len(t, resp2.Positions, 1)
		// Ensure we got different positions
		require.NotEqual(t, resp.Positions[0].Denom, resp2.Positions[0].Denom)
	})

	// Test invalid requests
	t.Run("nil request", func(t *testing.T) {
		_, err := queryServer.UserPosition(f.ctx, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid request")

		_, err = queryServer.UserPositions(f.ctx, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid request")
	})

	t.Run("empty user address", func(t *testing.T) {
		_, err := queryServer.UserPosition(f.ctx, &types.QueryUserPositionRequest{
			UserAddress: "",
			Denom:       "usdc",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "user address cannot be empty")

		_, err = queryServer.UserPositions(f.ctx, &types.QueryUserPositionsRequest{
			UserAddress: "",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "user address cannot be empty")
	})

	t.Run("empty denom", func(t *testing.T) {
		_, err := queryServer.UserPosition(f.ctx, &types.QueryUserPositionRequest{
			UserAddress: user1,
			Denom:       "",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "denom cannot be empty")
	})
}