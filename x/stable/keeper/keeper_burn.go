// Package keeper provides the core logic for the stable module.
// This file implements the burning functionality.
package keeper

import (
	"context"
	"math"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/stable/types"
)

// Burn converts stablecoins back into a chosen collateral asset.
// It performs the following actions:
// 1. Validates that the sender is an authorized authority.
// 2. Validates that the requested output collateral denom is accepted by the module.
// 3. Transfers the stablecoins to be burned from the sender to the module account.
// 4. Burns the transferred stablecoins.
// 5. Calculates the amount of collateral to return based on the exponents of the stablecoin and collateral denoms.
//   - The exchange rate is determined by the difference in decimal places (exponents).
//
// 6. Transfers the calculated amount of collateral from the module account to the sender.
// It returns the total amount of collateral returned or an error if any step fails.
func (k Keeper) Burn(ctx context.Context, sender sdk.AccAddress, stableToBurn sdkmath.Int, outputDenom string) (sdk.Coins, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// 1. Validate sender is an authority
	if !k.isAuthority(sender, params) {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "sender %s is not an authority", sender.String())
	}

	// 2. Validate output denom
	if err := k.validateOutputDenom(outputDenom, params); err != nil {
		return nil, err
	}

	// 3. Transfer stable coins from sender to module
	burnCoins := sdk.NewCoins(sdk.NewCoin(params.StableDenom, stableToBurn))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, burnCoins); err != nil {
		return nil, err
	}

	// 4. Burn the stable coins
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins); err != nil {
		return nil, err
	}

	// 5. Calculate collateral to return
	stableExponent, err := k.getDenomExponent(ctx, params.StableDenom)
	if err != nil {
		return nil, err
	}
	collateralExponent, err := k.getDenomExponent(ctx, outputDenom)
	if err != nil {
		return nil, err
	}

	// collateralToReturn = stableToBurn * 10^(collateralExponent - stableExponent)
	power := int(collateralExponent) - int(stableExponent)
	var collateralAmount sdkmath.Int
	if power == 0 {
		collateralAmount = stableToBurn
	} else if power > 0 {
		multiplier := sdkmath.NewIntFromUint64(uint64(math.Pow(10, float64(power))))
		collateralAmount = stableToBurn.Mul(multiplier)
	} else { // power < 0
		divisor := sdkmath.NewIntFromUint64(uint64(math.Pow(10, float64(abs(power)))))
		collateralAmount = stableToBurn.Quo(divisor)
	}

	if !collateralAmount.IsPositive() {
		return nil, errors.Wrapf(types.ErrInvalidAmount, "collateral to return must be positive, got %s", collateralAmount.String())
	}

	// 6. Transfer collateral from module to sender
	returnCoins := sdk.NewCoins(sdk.NewCoin(outputDenom, collateralAmount))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, returnCoins); err != nil {
		return nil, err
	}

	return returnCoins, nil
}

func (k Keeper) validateOutputDenom(outputDenom string, params types.Params) error {
	acceptedDenoms := make(map[string]struct{})
	for _, denom := range params.AcceptedDenoms {
		acceptedDenoms[denom] = struct{}{}
	}

	if _, ok := acceptedDenoms[outputDenom]; !ok {
		return errors.Wrapf(types.ErrInvalidDenom, "denom %s is not an accepted collateral for burning", outputDenom)
	}

	return nil
}
