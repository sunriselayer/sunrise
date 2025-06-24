package ante

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	feekeeper "github.com/sunriselayer/sunrise/x/fee/keeper"
)

type FeeDenomValidationDecorator struct {
	feeKeeper *feekeeper.Keeper
}

// Validates len(fee) == 1 and fee[0].Denom == feeKeeper.Params.FeeDenom
func NewFeeDenomValidationDecorator(feeKeeper *feekeeper.Keeper) FeeDenomValidationDecorator {
	return FeeDenomValidationDecorator{
		feeKeeper: feeKeeper,
	}
}

func (fdvd FeeDenomValidationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	fee := feeTx.GetFee()

	if fee.Empty() {
		return next(ctx, tx, simulate)
	}

	// check len(fee) == 1 and fee[0].Denom == feeKeeper.Params.FeeDenom
	if len(fee) > 1 {
		return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "only one fee denomination is allowed")
	}
	params, err := fdvd.feeKeeper.Params.Get(ctx)
	if err != nil {
		return ctx, err
	}
	if fee[0].Denom != params.FeeDenom {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "fee denom must be %s: %s", params.FeeDenom, fee[0].Denom)
	}

	return next(ctx, tx, simulate)
}
