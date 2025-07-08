// Package keeper provides the core logic for the stable module.
// This file implements the minting functionality.
package keeper

import (
	"context"
	"slices"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/stable/types"
)

// Mint creates new stablecoins and sends them to the sender.
// This operation can only be performed by an authority address.
// The authority is responsible for managing any collateral backing the minted stablecoins externally.
// It performs the following actions:
// 1. Validates that the sender is an authority address.
// 2. Mints the specified amount of stablecoins to the module account.
// 3. Sends the newly minted stablecoins from the module account to the sender.
// It returns the total amount of stablecoins minted or an error if any step fails.
func (k Keeper) Mint(ctx context.Context, sender sdk.AccAddress, amountToMint sdkmath.Int) (sdk.Coins, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// 1. Validates that the sender is an authority address.
	if !k.isAuthority(sender, params) {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "sender %s is not an authority address", sender.String())
	}

	if !amountToMint.IsPositive() {
		return nil, errors.Wrap(types.ErrInvalidAmount, "mint amount must be positive")
	}

	// 2. Mint stable coins to module
	stableCoins := sdk.NewCoins(sdk.NewCoin(params.StableDenom, amountToMint))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, stableCoins); err != nil {
		return nil, err
	}

	// 3. Send minted coins to sender
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, stableCoins); err != nil {
		return nil, err
	}

	return stableCoins, nil
}

func (k Keeper) isAuthority(sender sdk.AccAddress, params types.Params) bool {
	return slices.Contains(params.AuthorityAddresses, sender.String())
}
