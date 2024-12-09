package keeper

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/vmint/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	return nil
}
