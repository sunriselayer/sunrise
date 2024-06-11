package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/testutil/sample"
)

func TestMsgCollectIncentiveRewards_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCollectIncentiveRewards
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCollectIncentiveRewards{
				Sender: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCollectIncentiveRewards{
				Sender: sample.AccAddress(),
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
