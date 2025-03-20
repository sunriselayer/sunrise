package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k Keeper) WithdrawUnbonded(ctx sdk.Context, unstaking types.Unstaking) error {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return nil
	}

	// Convert bond denom to fee denom
	err = k.tokenConverterKeeper.Convert(ctx, unstaking.Amount.Amount, moduleAddr)
	if err != nil {
		return err
	}

	feeCoin := sdk.NewCoin(params.FeeDenom, unstaking.Amount.Amount)

	// Send coin to sender
	sender, err := k.addressCodec.StringToBytes(unstaking.Address)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(feeCoin))
	if err != nil {
		return err
	}

	return nil
}
