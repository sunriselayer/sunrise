package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k Keeper) BeginBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	// Handle module account rewards
	cacheCtx, write := sdk.UnwrapSDKContext(ctx).CacheContext()
	err := k.HandleModuleAccountRewards(cacheCtx)
	if err != nil {
		k.Logger().Error("failed to handle module account rewards", "error", err)
		return err
	}

	write()
	return nil
}

func (k Keeper) EndBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	// Withdraw unbonded
	cacheCtx, write := sdk.UnwrapSDKContext(ctx).CacheContext()
	err := k.GarbageCollectUnbonded(cacheCtx)
	if err != nil {
		k.Logger().Error("failed to garbage collect unbonded", "error", err)
		return err
	}

	write()
	return nil
}
