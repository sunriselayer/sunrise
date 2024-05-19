package keeper

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k Keeper) CalculateAmountOutFeeToken(ctx context.Context, amountInGovToken math.Int) math.Int {
	params := k.GetParams(ctx)

	supplyFee := k.bankKeeper.GetSupply(ctx, params.DenomFeeToken)
	supplyGov := k.bankKeeper.GetSupply(ctx, params.DenomGovToken)

	output := types.CalculateAmountOutFeeToken(supplyFee.Amount, supplyGov.Amount, amountInGovToken)

	return output
}

func (k Keeper) CalculateAmountInGovToken(ctx context.Context, amountOutFeeToken math.Int) math.Int {
	params := k.GetParams(ctx)

	supplyFee := k.bankKeeper.GetSupply(ctx, params.DenomFeeToken)
	supplyGov := k.bankKeeper.GetSupply(ctx, params.DenomGovToken)

	input := types.CalculateAmountInGovToken(supplyFee.Amount, supplyGov.Amount, amountOutFeeToken)

	return input
}

func (k Keeper) Convert(ctx context.Context, amountInGovToken math.Int, amountOutFeeToken math.Int, address sdk.AccAddress) error {
	params := k.GetParams(ctx)

	govToken := sdk.NewCoin(params.DenomGovToken, amountInGovToken)
	if err := govToken.Validate(); err != nil {
		return err
	}
	feeToken := sdk.NewCoin(params.DenomFeeToken, amountOutFeeToken)
	if err := feeToken.Validate(); err != nil {
		return err
	}

	supplyFee := k.bankKeeper.GetSupply(ctx, params.DenomFeeToken)

	if supplyFee.Amount.Add(amountOutFeeToken).GT(params.SupplyCapFeeToken) {
		return types.ErrExceedsSupplyCap
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(govToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(govToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(feeToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(feeToken)); err != nil {
		return err
	}

	return nil
}
