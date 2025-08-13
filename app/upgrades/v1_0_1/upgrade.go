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
	"github.com/sunriselayer/sunrise/x/lockup/keeper"

	"cosmossdk.io/collections"
)

// CreateUpgradeHandler creates an upgrade handler for the v1.0.1 upgrade.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bankKeeper bankkeeper.Keeper,
	lockupKeeper keeper.Keeper,
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

		// Transfer lockup account ownership
		oldLockupId := uint64(2)
		oldLockupAccount, err := lockupKeeper.GetLockupAccount(ctx, senderAddress, oldLockupId)
		if err != nil {
			return nil, fmt.Errorf("failed to get lockup account: %w", err)
		}

		// Get the address of the old lockup module account
		oldLockupModuleAddress, err := sdk.AccAddressFromBech32(oldLockupAccount.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to parse old lockup module address: %w", err)
		}

		// Get all funds from the old lockup module account
		balances := bankKeeper.GetAllBalances(ctx, oldLockupModuleAddress)

		// Get a new lockup ID for the new owner
		newLockupId, _, err := lockupKeeper.GetAndIncrementNextLockupAccountID(ctx, recipientAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to get next lockup account ID for new owner: %w", err)
		}

		// Generate the address for the new lockup module account
		newLockupModuleAddress := lockupKeeper.LockupAccountAddress(recipientAddress, newLockupId)

		// Transfer the funds from the old module account to the new module account
		if !balances.IsZero() {
			if err := bankKeeper.SendCoins(ctx, oldLockupModuleAddress, newLockupModuleAddress, balances); err != nil {
				return nil, fmt.Errorf("failed to send coins from old to new lockup account: %w", err)
			}
		}

		// Remove the old lockup account
		if err := lockupKeeper.LockupAccounts.Remove(ctx, collections.Join(senderAddress.Bytes(), oldLockupAccount.Id)); err != nil {
			return nil, fmt.Errorf("failed to remove old lockup account: %w", err)
		}

		// Create and set the new lockup account
		newLockupAccount := oldLockupAccount
		newLockupAccount.Owner = recipientAddress.String()
		newLockupAccount.Id = newLockupId
		newLockupAccount.Address = newLockupModuleAddress.String()

		if err := lockupKeeper.SetLockupAccount(ctx, newLockupAccount); err != nil {
			return nil, fmt.Errorf("failed to set new lockup account: %w", err)
		}

		ctx.Logger().Info(fmt.Sprintf("successfully transferred lockup account %d from %s to %s (new id: %d)", oldLockupId, senderAddress.String(), recipientAddress.String(), newLockupId))

		return mm.RunMigrations(goCtx, configurator, fromVM)
	}
}
