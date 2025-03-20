package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k Keeper) EndBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	// TODO: HandleModuleAccountRewards

	// Withdraw unbonded
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := k.IterateCompletedUnstakings(ctx, sdkCtx.BlockTime(), func(id uint64, value types.Unstaking) (stop bool, err error) {
		err = k.WithdrawUnbonded(ctx, value)
		if err != nil {
			return true, err
		}

		err = k.RemoveUnstaking(ctx, value.CompletionTime, id)
		if err != nil {
			return true, err
		}

		return false, nil
	})
	if err != nil {
		return err
	}

	return nil
}
