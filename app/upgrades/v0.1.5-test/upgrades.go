package v0_1_5_test

import (
	context "context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/sunriselayer/sunrise/app/keepers"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(context context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(context)
		ctx.Logger().Info(fmt.Sprintf("update start:%s", UpgradeName))

		err := upgradeSendCoin(ctx, keepers.BankKeeper)
		if err != nil {
			panic(err)
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
