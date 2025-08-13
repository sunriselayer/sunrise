// This file contains the upgrade handler for the v1.0.1 upgrade.
package v1_0_1

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// CreateUpgradeHandler creates an upgrade handler for the v1.0.1 upgrade.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bankKeeper bankkeeper.Keeper,
) types.UpgradeHandler {
	return func(goCtx context.Context, plan types.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(goCtx)
		ctx.Logger().Info(fmt.Sprintf("update start:%s", UpgradeName))

		senderAddressString := "sunrise1yfsg0ahx7dg99ytq4aqxvashmwq9tx3upp422w"
		recipientAddressString := "sunrise1pp2ruuhs0k7ayaxjupwj4k5qmgh0d72w8zu30p" // placeholder

		senderAddress, err := sdk.AccAddressFromBech32(senderAddressString)
		if err != nil {
			return nil, fmt.Errorf("error parsing sender address: %w", err)
		}

		recipientAddress, err := sdk.AccAddressFromBech32(recipientAddressString)
		if err != nil {
			return nil, fmt.Errorf("error parsing recipient address: %w", err)
		}

		coinsToSend := sdk.NewCoins(sdk.NewCoin("urise", math.NewInt(48_000_000_000)))

		if err := bankKeeper.SendCoins(ctx, senderAddress, recipientAddress, coinsToSend); err != nil {
			return nil, fmt.Errorf("failed to send coins: %w", err)
		}

		ctx.Logger().Info(fmt.Sprintf("successfully transferred %s from %s to %s", coinsToSend.String(), senderAddress.String(), recipientAddress.String()))

		return mm.RunMigrations(goCtx, configurator, fromVM)
	}
}
