package types

// import (
// 	"testing"

// 	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
// 	"github.com/stretchr/testify/require"
// 	"github.com/sunriselayer/sunrise/testutil/sample"
// )

// func TestMsgRegisterEvmAddress_ValidateBasic(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		msg  MsgRegisterEvmAddress
// 		err  error
// 	}{
// 		{
// 			name: "invalid address",
// 			msg:  MsgRegisterEvmAddress{},
// 			err:  sdkerrors.ErrInvalidAddress,
// 		}, {
// 			name: "valid address",
// 			msg:  MsgRegisterEvmAddress{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.msg.ValidateBasic()
// 			if tt.err != nil {
// 				require.ErrorIs(t, err, tt.err)
// 				return
// 			}
// 			require.NoError(t, err)
// 		})
// 	}
// }

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestValidateBasic(t *testing.T) {
	valAddr, err := sdk.ValAddressFromBech32("cosmosvaloper1xcy3els9ua75kdm783c3qu0rfa2eples6eavqq")
	require.NoError(t, err)
	evmAddr := common.BytesToAddress([]byte("hello"))

	msg := NewMsgRegisterEvmAddress(valAddr, evmAddr)
	require.NoError(t, msg.ValidateBasic())
	msg = &MsgRegisterEvmAddress{valAddr.String(), "invalid evm address"}
	require.Error(t, msg.ValidateBasic())
	msg = &MsgRegisterEvmAddress{"invalid validator address", evmAddr.Hex()}
	require.Error(t, msg.ValidateBasic())
}
