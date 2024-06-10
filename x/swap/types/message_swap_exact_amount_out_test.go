package types

import (
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/testutil/sample"
)

func TestMsgSwapExactAmountOut_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSwapExactAmountOut
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSwapExactAmountOut{
				Sender:            "invalid_address",
				InterfaceProvider: sample.AccAddress(),
				Route: Route{
					DenomIn:  "base",
					DenomOut: "quote",
					Strategy: &Route_Pool{
						Pool: &RoutePool{
							PoolId: 1,
						},
					},
				},
				MaxAmountIn: math.NewInt(1000000),
				AmountOut:   math.NewInt(1000000),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSwapExactAmountOut{
				Sender:            sample.AccAddress(),
				InterfaceProvider: sample.AccAddress(),
				Route: Route{
					DenomIn:  "base",
					DenomOut: "quote",
					Strategy: &Route_Pool{
						Pool: &RoutePool{
							PoolId: 1,
						},
					},
				},
				MaxAmountIn: math.NewInt(1000000),
				AmountOut:   math.NewInt(1000000),
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
