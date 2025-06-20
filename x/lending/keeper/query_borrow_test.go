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

func TestBorrowQuery(t *testing.T) {
	f := initFixture(t)
	defer f.ctrl.Finish()

	queryServer := keeper.NewQueryServerImpl(f.keeper)

	// Create test addresses
	borrower1 := sdk.AccAddress("borrower1").String()
	borrower2 := sdk.AccAddress("borrower2").String()

	// Create test borrows
	borrows := []types.Borrow{
		{
			Id:                   0,
			Borrower:             borrower1,
			Amount:               sdk.NewCoin("usdc", math.NewInt(100000)),
			CollateralPoolId:     1,
			CollateralPositionId: 100,
			BlockHeight:          100,
		},
		{
			Id:                   1,
			Borrower:             borrower1,
			Amount:               sdk.NewCoin("atom", math.NewInt(50000)),
			CollateralPoolId:     2,
			CollateralPositionId: 200,
			BlockHeight:          150,
		},
		{
			Id:                   2,
			Borrower:             borrower2,
			Amount:               sdk.NewCoin("usdc", math.NewInt(200000)),
			CollateralPoolId:     1,
			CollateralPositionId: 300,
			BlockHeight:          200,
		},
	}

	for _, borrow := range borrows {
		err := f.keeper.Borrows.Set(f.ctx, borrow.Id, borrow)
		require.NoError(t, err)
	}

	// Test Query single borrow
	t.Run("query single borrow", func(t *testing.T) {
		resp, err := queryServer.Borrow(f.ctx, &types.QueryBorrowRequest{
			BorrowId: 0,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, uint64(0), resp.Borrow.Id)
		require.Equal(t, borrower1, resp.Borrow.Borrower)
		require.Equal(t, "usdc", resp.Borrow.Amount.Denom)
		require.Equal(t, math.NewInt(100000), resp.Borrow.Amount.Amount)
		require.Equal(t, uint64(1), resp.Borrow.CollateralPoolId)
		require.Equal(t, uint64(100), resp.Borrow.CollateralPositionId)
		require.Equal(t, int64(100), resp.Borrow.BlockHeight)
	})

	// Test Query non-existent borrow
	t.Run("query non-existent borrow", func(t *testing.T) {
		_, err := queryServer.Borrow(f.ctx, &types.QueryBorrowRequest{
			BorrowId: 999,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "borrow not found")
	})

	// Test Query all borrows for a user
	t.Run("query all borrows for borrower1", func(t *testing.T) {
		resp, err := queryServer.UserBorrows(f.ctx, &types.QueryUserBorrowsRequest{
			Borrower: borrower1,
			Pagination: &query.PageRequest{
				Limit: 10,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Borrows, 2)
		
		// Check that we only got borrower1's borrows
		for _, borrow := range resp.Borrows {
			require.Equal(t, borrower1, borrow.Borrower)
		}
		
		// Check IDs
		ids := []uint64{resp.Borrows[0].Id, resp.Borrows[1].Id}
		require.Contains(t, ids, uint64(0))
		require.Contains(t, ids, uint64(1))
	})

	// Test Query all borrows for borrower2
	t.Run("query all borrows for borrower2", func(t *testing.T) {
		resp, err := queryServer.UserBorrows(f.ctx, &types.QueryUserBorrowsRequest{
			Borrower: borrower2,
			Pagination: &query.PageRequest{
				Limit: 10,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Borrows, 1)
		require.Equal(t, borrower2, resp.Borrows[0].Borrower)
		require.Equal(t, uint64(2), resp.Borrows[0].Id)
	})

	// Test pagination
	t.Run("query borrows with pagination", func(t *testing.T) {
		// Query first page with limit 1
		resp, err := queryServer.UserBorrows(f.ctx, &types.QueryUserBorrowsRequest{
			Borrower: borrower1,
			Pagination: &query.PageRequest{
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Borrows, 1)
		require.NotNil(t, resp.Pagination)
		require.NotEmpty(t, resp.Pagination.NextKey)

		// Query second page
		resp2, err := queryServer.UserBorrows(f.ctx, &types.QueryUserBorrowsRequest{
			Borrower: borrower1,
			Pagination: &query.PageRequest{
				Key:   resp.Pagination.NextKey,
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp2)
		require.Len(t, resp2.Borrows, 1)
		// Ensure we got different borrows
		require.NotEqual(t, resp.Borrows[0].Id, resp2.Borrows[0].Id)
	})

	// Test invalid requests
	t.Run("nil request", func(t *testing.T) {
		_, err := queryServer.Borrow(f.ctx, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid request")

		_, err = queryServer.UserBorrows(f.ctx, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid request")
	})

	t.Run("empty borrower address", func(t *testing.T) {
		_, err := queryServer.UserBorrows(f.ctx, &types.QueryUserBorrowsRequest{
			Borrower: "",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "borrower address cannot be empty")
	})
}