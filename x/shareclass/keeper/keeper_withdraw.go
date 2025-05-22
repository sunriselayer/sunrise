package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k Keeper) GarbageCollectUnbonded(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := k.IterateCompletedUnbondings(ctx, sdkCtx.BlockTime(), func(id uint64, value types.Unbonding) (stop bool, err error) {
		err = k.WithdrawUnbonded(ctx, value)
		if err != nil {
			return true, err
		}

		err = k.RemoveUnbonding(ctx, id)
		if err != nil {
			return true, err
		}

		return false, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) WithdrawUnbonded(ctx context.Context, unbonding types.Unbonding) error {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return err
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return err
	}

	if unbonding.Amount.Denom != bondDenom {
		return types.ErrInvalidUnbondedDenom
	}

	// Convert bond denom to fee denom
	err = k.tokenConverterKeeper.Convert(ctx, unbonding.Amount.Amount, moduleAddr)
	if err != nil {
		return err
	}

	feeCoin := sdk.NewCoin(feeDenom, unbonding.Amount.Amount)

	// Send coin to recipient
	recipient, err := k.addressCodec.StringToBytes(unbonding.RecipientAddress)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, sdk.NewCoins(feeCoin))
	if err != nil {
		return err
	}

	return nil
}
