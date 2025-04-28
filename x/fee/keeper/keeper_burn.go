package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sunriselayer/sunrise/x/fee/types"
)

func (k Keeper) Burn(ctx sdk.Context, fees sdk.Coins) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	burnRatio, err := math.LegacyNewDecFromStr(params.BurnRatio)
	if err != nil {
		return err
	}

	for _, fee := range fees {
		// skip if fee is not the fee denom
		if fee.Denom != params.FeeDenom {
			continue
		}
		burnAmount := burnRatio.MulInt(fee.Amount).TruncateInt()

		// skip if burn amount is zero
		if burnAmount.IsZero() {
			continue
		}

		burnCoin := sdk.NewCoin(fee.Denom, burnAmount)
		burnCoins := sdk.NewCoins(burnCoin)

		// burn coins from the fee module account
		if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx,
			authtypes.FeeCollectorName,
			types.ModuleName,
			burnCoins,
		); err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, err.Error())
		}

		// Event is emitted in the bank keeper
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins); err != nil {
			return err
		}
	}

	return nil
}
