package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k Keeper) EndBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	// Withdraw unbonded
	err := k.GarbageCollectUnbonded(ctx)
	if err != nil {
		return err
	}

	return nil
}
