// Package keeper provides the core logic for the stable module.
// This file implements the burning functionality.
package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/stable/types"
)

// Burn removes coins from circulation.
// This operation can only be performed by an authority address.
// It performs the following actions:
// 1. Validates that the sender is an authority address.
// 2. Transfers the coins to be burned from the sender to the module account.
// 3. Burns the transferred coins.
func (k Keeper) Burn(ctx context.Context, sender sdk.AccAddress, amountToBurn sdk.Coins) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// 1. Validate sender is an authority
	if !k.isAuthority(sender, params) {
		return errors.Wrapf(types.ErrInvalidSigner, "sender %s is not an authority", sender.String())
	}

	if !amountToBurn.IsAllPositive() {
		return errors.Wrap(types.ErrInvalidAmount, "burn amount must be positive")
	}

	// 2. Transfer coins from sender to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amountToBurn); err != nil {
		return err
	}

	// 3. Burn the coins
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, amountToBurn); err != nil {
		return err
	}

	return nil
}
