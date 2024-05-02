package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

func (k Keeper) TransferFromAccountToPoolModuleAccount(ctx context.Context, token sdk.Coin, address sdk.AccAddress, poolId uint64) error {
	moduleName := types.PoolModuleName(poolId)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, moduleName, sdk.NewCoins(token))

	return err
}

func (k Keeper) TransferFromPoolModuleAccountToAccount(ctx context.Context, token sdk.Coin, address sdk.AccAddress, poolId uint64) error {
	moduleName := types.PoolModuleName(poolId)
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleName, address, sdk.NewCoins(token))

	return err
}

func (k Keeper) TransferFromPoolModuleAccountToPoolTreasuryModuleAccount(ctx context.Context, token sdk.Coin, poolId uint64) error {
	moduleName := types.PoolModuleName(poolId)
	treasuryModuleName := types.PoolTreasuryModuleName(poolId)
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, moduleName, treasuryModuleName, sdk.NewCoins(token))

	return err
}

func (k Keeper) SwapExactAmountInSinglePool(ctx context.Context, poolId uint64, tokenIn sdk.Coin, denomOutConfirmation string, address sdk.AccAddress) (*math.Int, error) {
	// equal to zero is not caught here
	if tokenIn.Amount.LT(math.ZeroInt()) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid tokenIn amount %s", tokenIn.Amount.String())
	}

	// get pool
	pool, found := k.GetPool(ctx, poolId)

	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "pool id %d not found", poolId)
	}

	// check denom
	if tokenIn.Denom != pool.BaseDenom && tokenIn.Denom != pool.QuoteDenom {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid tokenIn denom %s", tokenIn.Denom)
	}

	// check denom
	var denomOut string
	if tokenIn.Denom == pool.BaseDenom {
		denomOut = pool.QuoteDenom
	} else {
		denomOut = pool.BaseDenom
	}

	if denomOut != denomOutConfirmation {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid denomOutConfirmation %s", denomOutConfirmation)
	}

	// pass zero input
	if tokenIn.Amount.Equal(math.ZeroInt()) {
		amountOut := math.ZeroInt()
		return &amountOut, nil
	}

	// transfer tokenIn from address to module
	if err := k.TransferFromAccountToPoolModuleAccount(ctx, tokenIn, address, poolId); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "error transferring tokenIn: %s", err.Error())
	}

	// calculate amount
	var amountOutNeg *math.Int
	var err error
	if tokenIn.Denom == pool.BaseDenom {
		amountOutNeg, err = types.CalculateDy(nil, tokenIn.Amount, nil, "")
	} else {
		amountOutNeg, err = types.CalculateDx(nil, tokenIn.Amount, nil, "")
	}

	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "error calculating amountOut: %s", err.Error())
	}

	// equal to zero is not caught here
	amountOut := amountOutNeg.Neg()
	if amountOut.LT(math.ZeroInt()) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amountOut is negative")
	}

	// calculate fees
	treasuryTaxAmount := k.GetParams(ctx).TreasuryTaxRate.MulInt(amountOut).RoundInt()
	poolFeeAmount := pool.FeeRate.MulInt(amountOut).RoundInt()
	amountOut = amountOut.Sub(treasuryTaxAmount).Sub(poolFeeAmount)

	// transfer tokenOut from module to address
	if amountOut.GT(math.ZeroInt()) {
		tokenOut := sdk.NewCoin(denomOut, amountOut)
		if err := k.TransferFromPoolModuleAccountToAccount(ctx, tokenOut, address, poolId); err != nil {
			return nil, err
		}
	}

	// transfer treasury tax to treasury
	if treasuryTaxAmount.GT(math.ZeroInt()) {
		treasuryTax := sdk.NewCoin(denomOut, treasuryTaxAmount)
		if err := k.TransferFromPoolModuleAccountToPoolTreasuryModuleAccount(ctx, treasuryTax, poolId); err != nil {
			return nil, err
		}
	}

	if poolFeeAmount.GT(math.ZeroInt()) {
		// TODO: emit event of pool fee
	}

	return &amountOut, nil
}

func (k Keeper) SwapExactAmountOutSinglePool(ctx context.Context, poolId uint64, tokenOut sdk.Coin, denomIntConfirmation string, address sdk.AccAddress) (*math.Int, error) {
	// equal to zero is caught here
	if tokenOut.Amount.LTE(math.ZeroInt()) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid tokenOut amount %s", tokenOut.Amount.String())
	}

	// get pool
	pool, found := k.GetPool(ctx, poolId)

	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "pool id %d not found", poolId)
	}

	// check denom
	if tokenOut.Denom != pool.BaseDenom && tokenOut.Denom != pool.QuoteDenom {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid tokenOut denom %s", tokenOut.Denom)
	}

	// check denom
	var denomIn string
	if tokenOut.Denom == pool.BaseDenom {
		denomIn = pool.QuoteDenom
	} else {
		denomIn = pool.BaseDenom
	}

	if denomIn != denomIntConfirmation {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid denomIntConfirmation %s", denomIntConfirmation)
	}

	// calculate amount
	treasuryTaxRate := k.GetParams(ctx).TreasuryTaxRate
	tokenOutBeforeFee := math.LegacyOneDec().Quo(math.LegacyOneDec().Sub(treasuryTaxRate).Sub(pool.FeeRate)).MulInt(tokenOut.Amount).RoundInt()
	tokenOutBeforeFeeNeg := tokenOutBeforeFee.Neg()

	var amountIn *math.Int
	var err error
	if tokenOut.Denom == pool.BaseDenom {
		amountIn, err = types.CalculateDy(nil, tokenOutBeforeFeeNeg, nil, "")
	} else {
		amountIn, err = types.CalculateDx(nil, tokenOutBeforeFeeNeg, nil, "")
	}

	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "error calculating amountOut: %s", err.Error())
	}

	// equal to zero is caught here
	if amountIn.LTE(math.ZeroInt()) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amountIn is negative")
	}

	// calculate fees
	treasuryTaxAmount := treasuryTaxRate.MulInt(tokenOutBeforeFee).RoundInt()
	poolFeeAmount := tokenOutBeforeFee.Sub(tokenOut.Amount).Sub(treasuryTaxAmount)

	// transfer tokenOut from module to address
	// no need to check zero case
	if err := k.TransferFromPoolModuleAccountToAccount(ctx, tokenOut, address, poolId); err != nil {
		return nil, err
	}

	// transfer treasury tax to treasury
	if treasuryTaxAmount.GT(math.ZeroInt()) {
		treasuryTax := sdk.NewCoin(tokenOut.Denom, treasuryTaxAmount)
		if err := k.TransferFromPoolModuleAccountToPoolTreasuryModuleAccount(ctx, treasuryTax, poolId); err != nil {
			return nil, err
		}
	}

	if poolFeeAmount.GT(math.ZeroInt()) {
		// TODO: emit event of pool fee
	}

	return amountIn, nil
}
