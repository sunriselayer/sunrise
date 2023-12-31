package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ethereum/go-ethereum/common"
)

var _ sdk.Msg = &MsgRegisterEvmAddress{}

func NewMsgRegisterEvmAddress(valAddress sdk.ValAddress, evmAddress common.Address) *MsgRegisterEvmAddress {
	return &MsgRegisterEvmAddress{
		Address:    sdk.AccAddress(valAddress).String(),
		EvmAddress: evmAddress.Hex(),
	}
}

func (msg *MsgRegisterEvmAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
