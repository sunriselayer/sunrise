package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) EndBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := k.DeleteExpiredBlobDeclarations(sdkCtx)
	if err != nil {
		k.Logger().Error("failed to delete expired blob declarations", "error", err)
	}

	// IF STATUS_VERIFIED is overtime, remove from store
	err = k.DeleteExpiredBlobIncludeds(sdkCtx)
	if err != nil {
		k.Logger().Error("failed to delete expired blob includeds", "error", err)
	}

	return nil
}

func (k Keeper) DeleteExpiredBlobDeclarations(ctx sdk.Context) error {

	return nil
}

func (k Keeper) DeleteExpiredBlobIncludeds(ctx sdk.Context) error {

	return nil
}
