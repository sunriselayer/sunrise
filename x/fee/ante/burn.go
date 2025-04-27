package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	feekeeper "github.com/sunriselayer/sunrise/x/fee/keeper"
)

type BurnFeeDecorator struct {
	feeKeeper *feekeeper.Keeper
}

func NewBurnFeeDecorator(feeKeeper *feekeeper.Keeper) BurnFeeDecorator {
	return BurnFeeDecorator{
		feeKeeper: feeKeeper,
	}
}

// CONTRACT: len(fee) == 1 and fee[0].Denom == params.FeeDenom
func (bfd BurnFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	fee := feeTx.GetFee()
	if !simulate {
		bfd.feeKeeper.Burn(ctx, fee)
	}

	return next(ctx, tx, simulate)
}
