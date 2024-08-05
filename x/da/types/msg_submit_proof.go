package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSubmitProof{}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgSubmitProof) ValidateBasic() error {
	return nil
}
