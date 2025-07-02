package keeper

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k Keeper) Convert(ctx context.Context, amount math.Int, address sdk.AccAddress) error {
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}
	fromDenom := params.FromDenom
	toDenom := params.ToDenom

	fromToken := sdk.NewCoin(fromDenom, amount)
	if err := fromToken.Validate(); err != nil {
		return err
	}
	toToken := sdk.NewCoin(toDenom, amount)
	if err := toToken.Validate(); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(fromToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(fromToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(toToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(toToken)); err != nil {
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
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}
	fromDenom := params.FromDenom
	toDenom := params.ToDenom

	fromToken := sdk.NewCoin(fromDenom, amount)
	if err := fromToken.Validate(); err != nil {
		return err
	}
	toToken := sdk.NewCoin(toDenom, amount)
	if err := toToken.Validate(); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(toToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(toToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(fromToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(fromToken)); err != nil {
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
