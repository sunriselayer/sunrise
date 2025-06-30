// Package keeper provides the core logic for the stable module.
// This file implements the minting functionality.
package keeper

import (
	"context"
	"math"
	"slices"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/stable/types"
)

// Mint converts collateral assets into stablecoins.
// It performs the following actions:
// 1. Validates that the sender is an authorized authority.
// 2. Validates that all collateral denoms are accepted by the module.
// 3. Transfers the collateral from the sender to the module account.
// 4. Calculates the amount of stablecoins to mint based on the exponents of the collateral and stablecoin denoms.
//   - The exchange rate is determined by the difference in decimal places (exponents).
//   - If metadata for a denom is not found, its exponent is treated as 0 (1:1 conversion against another denom with exponent 0).
//
// 5. Mints the calculated amount of stablecoins.
// 6. Sends the newly minted stablecoins to the sender.
// It returns the total amount of stablecoins minted or an error if any step fails.
func (k Keeper) Mint(ctx context.Context, sender sdk.AccAddress, collateral sdk.Coins) (sdk.Coins, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// 1. Validate sender is an authority
	if !k.isAuthority(sender, params) {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "sender %s is not an authority", sender.String())
	}

	// 2. Validate collateral denoms
	if err := k.validateCollateralDenoms(collateral, params); err != nil {
		return nil, err
	}

	// 3. Transfer collateral from sender to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, collateral); err != nil {
		return nil, err
	}

	// 4. Calculate amount to mint
	totalMintAmount := sdkmath.ZeroInt()
	stableExponent, err := k.getDenomExponent(ctx, params.StableDenom)
	if err != nil {
		return nil, err
	}

	for _, coin := range collateral {
		collateralExponent, err := k.getDenomExponent(ctx, coin.Denom)
		if err != nil {
			return nil, err
		}

		// amountToMint = collateralAmount * 10^(stableExponent - collateralExponent)
		power := int(stableExponent) - int(collateralExponent)

		var amountToMint sdkmath.Int
		if power == 0 {
			amountToMint = coin.Amount
		} else if power > 0 {
			multiplier := sdkmath.NewIntFromUint64(uint64(math.Pow(10, float64(power))))
			amountToMint = coin.Amount.Mul(multiplier)
		} else { // power < 0
			divisor := sdkmath.NewIntFromUint64(uint64(math.Pow(10, float64(abs(power)))))
			amountToMint = coin.Amount.Quo(divisor)
		}

		totalMintAmount = totalMintAmount.Add(amountToMint)
	}

	if !totalMintAmount.IsPositive() {
		return nil, errors.Wrap(types.ErrInvalidAmount, "mint amount must be positive")
	}

	// 5. Mint stable coins to module
	stableCoins := sdk.NewCoins(sdk.NewCoin(params.StableDenom, totalMintAmount))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, stableCoins); err != nil {
		return nil, err
	}

	// 6. Send minted coins to sender
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, stableCoins); err != nil {
		return nil, err
	}

	return stableCoins, nil
}

func (k Keeper) isAuthority(sender sdk.AccAddress, params types.Params) bool {
	return slices.Contains(params.AuthorityAddresses, sender.String())
}

func (k Keeper) validateCollateralDenoms(collateral sdk.Coins, params types.Params) error {
	acceptedDenoms := make(map[string]struct{})
	for _, denom := range params.AcceptedDenoms {
		acceptedDenoms[denom] = struct{}{}
	}

	if len(collateral) == 0 {
		return errors.Wrap(types.ErrInvalidDenom, "collateral cannot be empty")
	}

	for _, coin := range collateral {
		if _, ok := acceptedDenoms[coin.Denom]; !ok {
			return errors.Wrapf(types.ErrInvalidDenom, "denom %s is not an accepted collateral", coin.Denom)
		}
	}
	return nil
}

// getDenomExponent fetches the exponent for a given denom from the bank module's metadata.
// It returns the highest exponent found in the denom's units.
func (k Keeper) getDenomExponent(ctx context.Context, denom string) (uint32, error) {
	meta, found := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if !found {
		// As per spec, if metadata is not found, assume 1:1. This is handled by returning exponent 0.
		return 0, nil
	}

	var maxExponent uint32 = 0
	for _, unit := range meta.DenomUnits {
		if unit.Exponent > maxExponent {
			maxExponent = unit.Exponent
		}
	}
	return maxExponent, nil
}
