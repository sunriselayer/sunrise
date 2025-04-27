package ante

import (
	"bytes"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	feekeeper "github.com/sunriselayer/sunrise/x/fee/keeper"
	swapkeeper "github.com/sunriselayer/sunrise/x/swap/keeper"
	swaptypes "github.com/sunriselayer/sunrise/x/swap/types"
)

type SwapBeforeFeeDecorator struct {
	feeKeeper  *feekeeper.Keeper
	swapKeeper *swapkeeper.Keeper
}

func NewSwapBeforeFeeDecorator(feeKeeper *feekeeper.Keeper, swapKeeper *swapkeeper.Keeper) SwapBeforeFeeDecorator {
	return SwapBeforeFeeDecorator{
		feeKeeper:  feeKeeper,
		swapKeeper: swapKeeper,
	}
}

func (sbfd SwapBeforeFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	fee := feeTx.GetFee()

	// check len(fee) == 1 and fee[0].Denom == params.FeeDenom
	if len(fee) != 1 {
		return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "only one fee denomination is allowed")
	}
	params, err := sbfd.feeKeeper.Params.Get(ctx)
	if err != nil {
		return ctx, err
	}
	if fee[0].Denom != params.FeeDenom {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "fee denom must be %s: %s", params.FeeDenom, fee[0].Denom)
	}

	if !simulate {
		txConcrete, ok := tx.(*sdktx.Tx)
		if !ok {
			return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a sdktx.Tx")
		}

		for _, ext := range txConcrete.Body.ExtensionOptions {
			var swapExt swaptypes.SwapBeforeFeeExtension
			todo := true

			if todo {
				_, _, err = sbfd.swapKeeper.SwapExactAmountOut(
					ctx,
					sdk.AccAddress(feeTx.FeePayer()),
					authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
					swapExt.Route,
					swapExt.MaxAmountIn,
					fee[0].Amount,
				)
				if err != nil {
					return ctx, errorsmod.Wrapf(err, "failed to swap before fee")
				}

				break
			}
		}
	}

	return next(ctx, tx, simulate)
}
