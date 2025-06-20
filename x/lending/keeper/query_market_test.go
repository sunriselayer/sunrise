package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/lending/keeper"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func TestMarketQuery(t *testing.T) {
	f := initFixture(t)
	defer f.ctrl.Finish()

	queryServer := keeper.NewQueryServerImpl(f.keeper)

	// Create test markets
	markets := []types.Market{
		{
			Denom:             "usdc",
			TotalSupplied:     math.NewInt(1000000),
			TotalBorrowed:     math.NewInt(500000),
			GlobalRewardIndex: math.LegacyOneDec(),
			RiseDenom:         "riseusdc",
		},
		{
			Denom:             "atom",
			TotalSupplied:     math.NewInt(2000000),
			TotalBorrowed:     math.NewInt(1000000),
			GlobalRewardIndex: math.LegacyNewDecWithPrec(11, 1), // 1.1
			RiseDenom:         "riseatom",
		},
	}

	for _, market := range markets {
		err := f.keeper.SetMarket(f.ctx, market)
		require.NoError(t, err)
	}

	// Test Query single market
	t.Run("query single market", func(t *testing.T) {
		resp, err := queryServer.Market(f.ctx, &types.QueryMarketRequest{
			Denom: "usdc",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "usdc", resp.Market.Denom)
		require.Equal(t, math.NewInt(1000000), resp.Market.TotalSupplied)
		require.Equal(t, math.NewInt(500000), resp.Market.TotalBorrowed)
		require.Equal(t, "riseusdc", resp.Market.RiseDenom)
	})

	// Test Query non-existent market
	t.Run("query non-existent market", func(t *testing.T) {
		_, err := queryServer.Market(f.ctx, &types.QueryMarketRequest{
			Denom: "nonexistent",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "market not found")
	})

	// Test Query all markets
	t.Run("query all markets", func(t *testing.T) {
		resp, err := queryServer.Markets(f.ctx, &types.QueryMarketsRequest{
			Pagination: &query.PageRequest{
				Limit: 10,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Markets, 2)
		// Markets may be returned in any order
		denoms := []string{resp.Markets[0].Denom, resp.Markets[1].Denom}
		require.Contains(t, denoms, "usdc")
		require.Contains(t, denoms, "atom")
	})

	// Test pagination
	t.Run("query markets with pagination", func(t *testing.T) {
		// Query first page with limit 1
		resp, err := queryServer.Markets(f.ctx, &types.QueryMarketsRequest{
			Pagination: &query.PageRequest{
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Markets, 1)
		require.NotNil(t, resp.Pagination)
		require.NotEmpty(t, resp.Pagination.NextKey)

		// Query second page
		resp2, err := queryServer.Markets(f.ctx, &types.QueryMarketsRequest{
			Pagination: &query.PageRequest{
				Key:   resp.Pagination.NextKey,
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp2)
		require.Len(t, resp2.Markets, 1)
		// Ensure we got different markets
		require.NotEqual(t, resp.Markets[0].Denom, resp2.Markets[0].Denom)
	})

	// Test invalid requests
	t.Run("nil request", func(t *testing.T) {
		_, err := queryServer.Market(f.ctx, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid request")

		_, err = queryServer.Markets(f.ctx, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid request")
	})

	t.Run("empty denom", func(t *testing.T) {
		_, err := queryServer.Market(f.ctx, &types.QueryMarketRequest{
			Denom: "",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "denom cannot be empty")
	})
}