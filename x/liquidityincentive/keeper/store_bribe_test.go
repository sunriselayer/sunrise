package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNBribe(keeper keeper.Keeper, ctx context.Context, n int) []types.Bribe {
	// Ensure epoch 1 exists
	epoch := types.Epoch{
		Id:         1,
		StartBlock: 0,
		EndBlock:   100,
		Gauges:     []types.Gauge{},
	}
	_ = keeper.SetEpoch(ctx, epoch)

	items := make([]types.Bribe, n)
	for i := range items {
		items[i].EpochId = 1
		items[i].PoolId = uint64(i)
		items[i].Amount = sdk.NewCoins(sdk.NewCoin("test", math.OneInt()))
		items[i].ClaimedAmount = sdk.NewCoins(sdk.NewCoin("test", math.ZeroInt()))

		id, _ := keeper.AppendBribe(ctx, items[i])
		items[i].Id = id
	}
	return items
}

func TestBribeGet(t *testing.T) {
	f := initFixture(t)
	items := createNBribe(f.keeper, f.ctx, 10)
	for _, item := range items {
		rst, found, err := f.keeper.GetBribe(f.ctx, item.Id)
		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestBribeRemove(t *testing.T) {
	f := initFixture(t)
	items := createNBribe(f.keeper, f.ctx, 10)
	for _, item := range items {
		err := f.keeper.RemoveBribe(f.ctx, item.Id)
		require.NoError(t, err)
		_, found, err := f.keeper.GetBribe(f.ctx, item.PoolId)
		require.NoError(t, err)
		require.False(t, found)
	}
}

func TestBribeGetAll(t *testing.T) {
	f := initFixture(t)
	items := createNBribe(f.keeper, f.ctx, 10)
	bribes, err := f.keeper.GetAllBribes(f.ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(bribes),
	)
}

func TestBribeGetAllByEpochId(t *testing.T) {
	f := initFixture(t)
	items := createNBribe(f.keeper, f.ctx, 10)
	bribes, err := f.keeper.GetAllBribeByEpochId(f.ctx, 1)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(bribes),
	)
}
