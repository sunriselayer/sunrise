package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/testutil/sample"
)

func TestMsgCreateRegistration_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateRegistration
		err  error
	}{
		{
			name: "invalid liquidity provider address",
			msg: MsgCreateRegistration{
				LiquidityProvider: "invalid_address",
				Grantee:           sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid grantee address",
			msg: MsgCreateRegistration{
				LiquidityProvider: sample.AccAddress(),
				Grantee:           "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid liquidity provider address",
			msg: MsgCreateRegistration{
				LiquidityProvider: sample.AccAddress(),
				Grantee:           sample.AccAddress(),
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
			name: "invalid liquidity provider address",
			msg: MsgUpdateRegistration{
				LiquidityProvider: "invalid_address",
				Grantee:           sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid grantee address",
			msg: MsgUpdateRegistration{
				LiquidityProvider: sample.AccAddress(),
				Grantee:           "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid liquidity provider address",
			msg: MsgUpdateRegistration{
				LiquidityProvider: sample.AccAddress(),
				Grantee:           sample.AccAddress(),
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
				LiquidityProvider: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteRegistration{
				LiquidityProvider: sample.AccAddress(),
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
