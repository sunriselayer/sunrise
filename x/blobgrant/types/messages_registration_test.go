package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise-app/testutil/sample"
)

func TestMsgCreateRegistration_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateRegistration
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateRegistration{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateRegistration{
				Creator: sample.AccAddress(),
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

func TestMsgUpdateRegistration_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateRegistration
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateRegistration{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateRegistration{
				Creator: sample.AccAddress(),
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

func TestMsgDeleteRegistration_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteRegistration
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteRegistration{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteRegistration{
				Creator: sample.AccAddress(),
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
