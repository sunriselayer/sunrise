package keeper

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k Keeper) GetDenoms(ctx context.Context) (bondDenom string, feeDenom string, err error) {
	bondDenom, err = k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return "", "", err
	}
	feeDenom, err = k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return "", "", err
	}

	return bondDenom, feeDenom, nil
}

func (k Keeper) Convert(ctx context.Context, amount math.Int, address sdk.AccAddress) error {
	bondDenom, feeDenom, err := k.GetDenoms(ctx)
	if err != nil {
		return err
	}

	bondToken := sdk.NewCoin(bondDenom, amount)
	if err := bondToken.Validate(); err != nil {
		return err
	}
	feeToken := sdk.NewCoin(feeDenom, amount)
	if err := feeToken.Validate(); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(bondToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), sdk.NewCoins(bondToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(feeToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(feeToken)); err != nil {
		return err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventConvert{
		Address: address.String(),
		Amount:  amount.String(),
	}); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ConvertReverse(ctx context.Context, amount math.Int, address sdk.AccAddress) error {
	bondDenom, feeDenom, err := k.GetDenoms(ctx)
	if err != nil {
		return err
	}

	bondToken := sdk.NewCoin(bondDenom, amount)
	if err := bondToken.Validate(); err != nil {
		return err
	}
	feeToken := sdk.NewCoin(feeDenom, amount)
	if err := feeToken.Validate(); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(feeToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), sdk.NewCoins(feeToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(bondToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(bondToken)); err != nil {
		return err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventConvertReverse{
		Address: address.String(),
		Amount:  amount.String(),
	}); err != nil {
		return err
	}

	return nil
}
