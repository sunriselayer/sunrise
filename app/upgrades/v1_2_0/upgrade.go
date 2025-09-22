// This file contains the upgrade handler for the v1.2.0 upgrade.
package v1_2_0

import (
	"context"
	"fmt"

	"cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// CreateUpgradeHandler creates an upgrade handler for the v1.2.0 upgrade.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) types.UpgradeHandler {
	return func(goCtx context.Context, plan types.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)
		ctx.Logger().Info(fmt.Sprintf("update start:%s", UpgradeName))

		return mm.RunMigrations(goCtx, configurator, fromVM)
	}
}
