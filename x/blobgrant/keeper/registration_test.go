package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/blobgrant/keeper"
	"github.com/sunriselayer/sunrise/x/blobgrant/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNRegistration(keeper keeper.Keeper, ctx context.Context, n int) []types.Registration {
	items := make([]types.Registration, n)
	for i := range items {
		items[i].Address = strconv.Itoa(i)

		keeper.SetRegistration(ctx, items[i])
	}
	return items
}

func TestRegistrationGet(t *testing.T) {
	keeper, ctx := keepertest.GrantKeeper(t)
	items := createNRegistration(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetRegistration(ctx,
			item.Address,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestRegistrationRemove(t *testing.T) {
	keeper, ctx := keepertest.GrantKeeper(t)
	items := createNRegistration(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveRegistration(ctx,
			item.Address,
		)
		_, found := keeper.GetRegistration(ctx,
			item.Address,
		)
		require.False(t, found)
	}
}

func TestRegistrationGetAll(t *testing.T) {
	keeper, ctx := keepertest.GrantKeeper(t)
	items := createNRegistration(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllRegistration(ctx)),
	)
}
