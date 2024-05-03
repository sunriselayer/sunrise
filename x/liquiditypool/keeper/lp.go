package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	types "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k Keeper) joinPool(ctx context.Context, poolId uint64, baseToken sdk.Coin, quoteToken sdk.Coin, dryRun bool, sender *string, minShareAmount *math.Int) (*math.Int, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("pool id %d doesn't exist", poolId))
	}

	if baseToken.Denom != pool.BaseDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("base denom %s is invalid", baseToken.Denom))
	}
	if quoteToken.Denom != pool.QuoteDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("quote denom %s is invalid", quoteToken.Denom))
	}

	x, y := k.GetPoolBalance(ctx, pool)
	price, err := types.CalculatePrice(x, y, pool)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error calculating price: %s", err.Error()))
	}
	value := types.LpTokenValueInQuoteUnit(x, y, *price)

	newX := x.Add(baseToken.Amount)
	newY := y.Add(quoteToken.Amount)
	newPrice, err := types.CalculatePrice(newX, newY, pool)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error calculating price: %s", err.Error()))
	}
	newValue := types.LpTokenValueInQuoteUnit(newX, newY, *newPrice)

	if newValue.LTE(value) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "value is not increased")
	}

	supply := k.GetLpTokenSupply(ctx, pool.Id)
	newSupplyAmount := newValue.Quo(value).MulInt(supply.Amount).RoundInt()
	additionalSupply := sdk.NewCoin(supply.Denom, newSupplyAmount.Sub(supply.Amount))

	if !dryRun {
		address := sdk.MustAccAddressFromBech32(*sender)

		if additionalSupply.Amount.LT(*minShareAmount) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "min share amount is not met")
		}

		if err := k.TransferFromAccountToPoolModule(ctx, baseToken, address, pool.Id); err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error transferring base token: %s", err.Error()))
		}
		if err := k.TransferFromAccountToPoolModule(ctx, quoteToken, address, pool.Id); err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error transferring quote token: %s", err.Error()))
		}

		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(additionalSupply)); err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error minting lp token: %s", err.Error()))
		}
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(additionalSupply)); err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error sending lp token: %s", err.Error()))
		}
	}

	return &additionalSupply.Amount, nil
}

func (k Keeper) exitPool(ctx context.Context, poolId uint64, shareAmount math.Int, dryRun bool, sender *string, minAmountBase *math.Int, minAmountQuote *math.Int) ([]sdk.Coin, error) {

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("pool id %d doesn't exist", poolId))
	}

	_ = pool

	if !dryRun {
		address := sdk.MustAccAddressFromBech32(*sender)
		_ = address
	}

	return nil, nil
}
