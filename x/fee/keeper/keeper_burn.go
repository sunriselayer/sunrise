package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sunriselayer/sunrise/x/fee/types"
)

// Burn handles burning a portion of transaction fees in a fault-tolerant way.
// If an internal error like a swap failure occurs, the error is logged and the burn
// is skipped, leaving the parent transaction unaffected. The burn itself is an
// atomic operation to prevent funds from getting stuck.
//
// `fees` represents the total transaction fees, not the amount to be burned.
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

	// We use a cache context to make the burn process atomic.
	cacheCtx, write := ctx.CacheContext()
	if err := k.burnCoin(cacheCtx, burnCoin, params); err != nil {
		k.Logger().Error("failed to burn fees", "err", err)
		// Do not write cache context to main state if burning fails.
		return nil
	}

	// Write cache context to main state only if burning is successful.
	write()
	return nil
}

// burnCoin performs the actual burning of a coin. It is designed to be called
// within a cached context to ensure atomicity.
func (k Keeper) burnCoin(ctx sdk.Context, coin sdk.Coin, params types.Params) error {
	coins := sdk.NewCoins(coin)
	// Send coins to be burned from the fee collector to the fee module.
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx,
		authtypes.FeeCollectorName,
		types.ModuleName,
		coins,
	); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	if params.FeeDenom == params.BurnDenom {
		// burn coins from the fee module account
		// Event is emitted in the bank keeper
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
			return errorsmod.Wrap(err, "failed to burn coins")
		}
	} else {
		// swap to burn denom and burn
		pool, found, err := k.liquidityPoolKeeper.GetPool(ctx, params.BurnPoolId)
		if err != nil {
			return errorsmod.Wrap(err, "failed to get pool for burning")
		}
		if !found {
			return errorsmod.Wrapf(types.ErrPoolNotFound, "pool not found for burning: %d", params.BurnPoolId)
		}

		// swap to burn denom
		swappedAmount, err := k.liquidityPoolKeeper.SwapExactAmountIn(ctx,
			authtypes.NewModuleAddress(types.ModuleName),
			pool,
			coin,
			params.BurnDenom,
			true,
		)
		if err != nil {
			return errorsmod.Wrap(err, "failed to swap to burn denom")
		}

		// burn swapped coins from the fee module account
		// Event is emitted in the bank keeper
		swappedCoin := sdk.NewCoin(params.BurnDenom, swappedAmount)
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(swappedCoin)); err != nil {
			return errorsmod.Wrap(err, "failed to burn coins after swap")
		}
	}

	return nil
}
