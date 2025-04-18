package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k Keeper) EndBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	// Handle module account rewards
	err := k.HandleModuleAccountRewards(ctx)
	if err != nil {
		return err
	}

	// Withdraw unbonded
	err = k.GarbageCollectUnbonded(ctx)
	if err != nil {
		return err
	}

	return nil
}
