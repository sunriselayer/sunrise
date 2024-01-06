package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgMintDerivative{}

func NewMsgMintDerivative(
	sender sdk.AccAddress,
	validator sdk.ValAddress,
	amount sdk.Coin,
) *MsgMintDerivative {
	return &MsgMintDerivative{
		Sender:    sender.String(),
		Validator: validator.String(),
		Amount:    amount,
	}
}

func (msg *MsgMintDerivative) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	_, err = sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if msg.Amount.IsNil() || !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "'%s'", msg.Amount)
	}

	return nil
}
