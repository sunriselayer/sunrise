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
				Sender:    "invalid_address",
				MinAmount: math.NewInt(1),
				MaxAmount: math.NewInt(2),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgConvert{
				Sender:    sample.AccAddress(),
				MinAmount: math.NewInt(1),
				MaxAmount: math.NewInt(2),
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
