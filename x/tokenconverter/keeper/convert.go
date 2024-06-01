package keeper

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k Keeper) CalculateConversionAmount(ctx context.Context, minAmountOutFeeToken math.Int, maxAmountInGovToken math.Int) (math.Int, error) {
	params := k.GetParams(ctx)

	supplyFee := k.bankKeeper.GetSupply(ctx, params.FeeDenom)

	space := params.MaxSupplyFee.Sub(supplyFee.Amount)
	if space.IsZero() || space.LT(minAmountOutFeeToken) {
		return math.ZeroInt(), types.ErrExceedsMaxSupply
	}

	var amount math.Int
	if space.LT(maxAmountInGovToken) {
		amount = space
	} else {
		amount = maxAmountInGovToken
	}

	return amount, nil
}

func (k Keeper) BurnAndMint(ctx context.Context, amount math.Int, address sdk.AccAddress) error {
	params := k.GetParams(ctx)

	govToken := sdk.NewCoin(params.GovDenom, amount)
	if err := govToken.Validate(); err != nil {
		return err
	}
	feeToken := sdk.NewCoin(params.FeeDenom, amount)
	if err := feeToken.Validate(); err != nil {
		return err
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
