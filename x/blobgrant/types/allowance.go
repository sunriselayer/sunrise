package types

import (
	"context"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"cosmossdk.io/x/feegrant"
	blobtypes "github.com/sunriselayer/sunrise/x/blob/types"
)

var _ feegrant.FeeAllowanceI = FeeAllowance{}

func NewFeeAllowance(gas math.Int, expiry time.Time) FeeAllowance {
	return FeeAllowance{
		Gas:    gas,
		Expiry: expiry,
	}
}

func (allowance FeeAllowance) Accept(ctx context.Context, fee sdk.Coins, msgs []sdk.Msg) (remove bool, err error) {
	if len(fee) == 0 {
		return false, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot accept txs with empty fee")
	}
	if len(fee) > 1 {
		return false, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot accept txs with multiple fee tokens")
	}
	if fee[0].Denom != GrantTokenDenom {
		return false, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid fee token %s", fee[0].Denom)
	}

	for _, msg := range msgs {
		_, ok := msg.(*blobtypes.MsgPayForBlobs)
		if !ok {
			return false, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid message type")
		}
	}

	return true, nil
}

func (allowance FeeAllowance) ValidateBasic() error {
	if allowance.Gas.IsNil() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "gas must be positive: nil")
	}
	if !allowance.Gas.IsPositive() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "gas must be positive: %s", allowance.Gas)
	}

	return nil
}

func (allowance FeeAllowance) ExpiresAt() (*time.Time, error) {
	return &allowance.Expiry, nil
}
