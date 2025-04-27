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

// CONTRACT: len(fee) == 1 and fee[0].Denom == feeKeeper.Params.FeeDenom
func NewBurnFeeDecorator(feeKeeper *feekeeper.Keeper) BurnFeeDecorator {
	return BurnFeeDecorator{
		feeKeeper: feeKeeper,
	}
}

func (bfd BurnFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	fee := feeTx.GetFee()
	if !simulate {
		err := bfd.feeKeeper.Burn(ctx, fee)
		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}
