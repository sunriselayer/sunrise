// This file contains the upgrade handler that allows stuck swap packets to complete.
package v1_4_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// CreateUpgradeHandler creates the handler for UpgradeName.
// Completing in-flight packets is gated on this upgrade's done height, so the
// handler itself only runs module migrations.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(goCtx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)
		ctx.Logger().Info(fmt.Sprintf("upgrade start: %s at height %d", UpgradeName, ctx.BlockHeight()))
		return mm.RunMigrations(goCtx, configurator, fromVM)
	}
}
