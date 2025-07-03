package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sunriselayer/sunrise/x/fee/types"
)

// fees means whole tx fees, not amount to burn
func (k Keeper) Burn(ctx sdk.Context, fees sdk.Coins) error {
	err := fees.Validate()
	if err != nil {
		return err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	burnRatio, err := math.LegacyNewDecFromStr(params.BurnRatio)
	if err != nil {
		return err
	}

	found, feeCoin := fees.Find(params.FeeDenom)
	if !found {
		return nil
	}

	burnAmount := burnRatio.MulInt(feeCoin.Amount).TruncateInt()

	// skip if burn amount is zero
	if burnAmount.IsZero() {
		return nil
	}

	burnCoin := sdk.NewCoin(feeCoin.Denom, burnAmount)
	burnCoins := sdk.NewCoins(burnCoin)

	// Send coins to be burned from the fee collector to the fee module.
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx,
		authtypes.FeeCollectorName,
		types.ModuleName,
		burnCoins,
	); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	if params.FeeDenom == params.BurnDenom {
		// burn coins from the fee module account
		// Event is emitted in the bank keeper
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins); err != nil {
			return err
		}
	} else {
		// swap to burn denom and burn
		pool, found, err := k.liquidityPoolKeeper.GetPool(ctx, params.BurnPoolId)
		if err != nil {
			k.Logger().Error("failed to get pool for burning", "pool_id", params.BurnPoolId, "err", err)
			return nil
		}
		if !found {
			k.Logger().Error("pool not found for burning", "pool_id", params.BurnPoolId)
			return nil
		}

		// swap to burn denom
		swappedAmount, err := k.liquidityPoolKeeper.SwapExactAmountIn(ctx,
			authtypes.NewModuleAddress(types.ModuleName),
			pool,
			burnCoin,
			params.BurnDenom,
			false,
		)
		if err != nil {
			k.Logger().Error("failed to swap to burn denom", "err", err)
			return nil
		}

		// burn swapped coins from the fee module account
		// Event is emitted in the bank keeper
		swappedCoin := sdk.NewCoin(params.BurnDenom, swappedAmount)
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(swappedCoin)); err != nil {
			k.Logger().Error("failed to burn coins after swap", "err", err)
			return nil
		}
	}

	return nil
}
