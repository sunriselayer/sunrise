package types

import (
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/testutil/sample"
)

func TestMsgSwapExactAmountIn_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSwapExactAmountIn
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSwapExactAmountIn{
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
				AmountIn:     math.NewInt(1000000),
				MinAmountOut: math.NewInt(1000000),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSwapExactAmountIn{
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
				AmountIn:     math.NewInt(1000000),
				MinAmountOut: math.NewInt(1000000),
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
