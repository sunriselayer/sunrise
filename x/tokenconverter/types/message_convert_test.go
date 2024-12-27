package types

import (
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/testutil/sample"
)

func TestMsgConvert_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgConvert
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgConvert{
				Sender: "invalid_address",
				Amount: math.NewInt(1),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgConvert{
				Sender: sample.AccAddress(),
				Amount: math.NewInt(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
