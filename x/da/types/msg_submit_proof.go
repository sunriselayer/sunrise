package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSubmitProof{}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgSubmitProof) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	// check number of proofs <> indices
	if len(msg.Indices) != len(msg.Proofs) {
		return ErrIndicesAndProofsMismatch
	}

	return nil
}
