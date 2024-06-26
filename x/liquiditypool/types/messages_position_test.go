package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/testutil/sample"
)

func TestMsgCreatePosition_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreatePosition
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreatePosition{
				Sender: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreatePosition{
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

func TestMsgIncreaseLiquidity_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgIncreaseLiquidity
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgIncreaseLiquidity{
				Sender: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgIncreaseLiquidity{
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

func TestMsgDecreaseLiquidity_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDecreaseLiquidity
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDecreaseLiquidity{
				Sender: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDecreaseLiquidity{
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
