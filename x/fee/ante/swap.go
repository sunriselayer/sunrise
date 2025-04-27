package ante

import (
	errorsmod "cosmossdk.io/errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	feekeeper "github.com/sunriselayer/sunrise/x/fee/keeper"
	swapkeeper "github.com/sunriselayer/sunrise/x/swap/keeper"
	swaptypes "github.com/sunriselayer/sunrise/x/swap/types"
)

type SwapBeforeFeeDecorator struct {
	feeKeeper  *feekeeper.Keeper
	swapKeeper *swapkeeper.Keeper
}

// CONTRACT: len(fee) == 1 and fee[0].Denom == feeKeeper.Params.FeeDenom
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

	if !simulate {
		hasExtOptsTx, ok := tx.(authante.HasExtensionOptionsTx)
		if !ok {
			return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a HasExtensionOptionsTx")
		}

		once := false

		for _, ext := range hasExtOptsTx.GetExtensionOptions() {
			if ext.TypeUrl == codectypes.MsgTypeURL(&swaptypes.SwapBeforeFeeExtension{}) {
				if once {
					return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "swap before fee extension can only be used once")
				}

				// Very important
				if feeTx.FeeGranter() != nil {
					return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "fee granter is not allowed with swap before fee")
				}

				// Unmarshal the extension
				var swapExt swaptypes.SwapBeforeFeeExtension
				err := swapExt.Unmarshal(ext.Value)
				if err != nil {
					return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to unmarshal swap before fee extension")
				}

				// Validate route
				params, err := sbfd.feeKeeper.Params.Get(ctx)
				if err != nil {
					return ctx, err
				}

				if swapExt.Route.DenomOut != params.FeeDenom {
					return ctx, errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "route denom out must be %s: %s", params.FeeDenom, swapExt.Route.DenomOut)
				}

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

				once = true
			}
		}
	}

	return next(ctx, tx, simulate)
}
