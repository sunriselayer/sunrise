package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/testutil/sample"
)

func TestMsgCollectFees_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgClaimRewards
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgClaimRewards{
				Sender:      "invalid_address",
				PositionIds: []uint64{1},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgClaimRewards{
				Sender:      sample.AccAddress(),
				PositionIds: []uint64{1},
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
