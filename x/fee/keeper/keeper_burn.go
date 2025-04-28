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
	err := fees.Validate()
	if err != nil {
		return err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	burnRatio, err := math.LegacyNewDecFromStr(params.BurnRatio)
	if err != nil {
		return err
	}

	found, feeCoin := fees.Find(params.FeeDenom)
	if found {
		burnAmount := burnRatio.MulInt(feeCoin.Amount).TruncateInt()

		// skip if burn amount is zero
		if burnAmount.IsZero() {
			return nil
		}

		burnCoin := sdk.NewCoin(feeCoin.Denom, burnAmount)
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
