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
	nonTransferableDenom := params.NonTransferableDenom
	transferableDenom := params.TransferableDenom

	nonTransferableToken := sdk.NewCoin(nonTransferableDenom, amount)
	if err := nonTransferableToken.Validate(); err != nil {
		return err
	}
	transferableToken := sdk.NewCoin(transferableDenom, amount)
	if err := transferableToken.Validate(); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(nonTransferableToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(nonTransferableToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(transferableToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(transferableToken)); err != nil {
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
	nonTransferableDenom := params.NonTransferableDenom
	transferableDenom := params.TransferableDenom

	nonTransferableToken := sdk.NewCoin(nonTransferableDenom, amount)
	if err := nonTransferableToken.Validate(); err != nil {
		return err
	}
	transferableToken := sdk.NewCoin(transferableDenom, amount)
	if err := transferableToken.Validate(); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(transferableToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(transferableToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(nonTransferableToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(nonTransferableToken)); err != nil {
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
