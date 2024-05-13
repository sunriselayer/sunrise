package types_test

import (
	"context"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	blobtypes "github.com/sunriselayer/sunrise/x/blob/types"
	"github.com/sunriselayer/sunrise/x/blobgrant/types"
)

func TestFeeAllowanceAccept(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		allowance types.FeeAllowance
		msg       sdk.Msg
		fee       sdk.Coins
		err       error
	}{
		{
			name:      "empty fee",
			allowance: types.NewFeeAllowance(math.OneInt(), now),
			msg:       &blobtypes.MsgPayForBlobs{},
			fee:       sdk.Coins{},
			err:       errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot accept txs with empty fee"),
		},
		{
			name:      "multiple fees",
			allowance: types.NewFeeAllowance(math.OneInt(), now),
			msg:       &blobtypes.MsgPayForBlobs{},
			fee:       sdk.Coins{sdk.NewCoin(types.GrantTokenDenom, math.OneInt()), sdk.NewCoin("zzz", math.OneInt())},
			err:       errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot accept txs with multiple fee tokens"),
		},
		{
			name:      "invalid fee denom",
			allowance: types.NewFeeAllowance(math.OneInt(), now),
			msg:       &blobtypes.MsgPayForBlobs{},
			fee:       sdk.Coins{sdk.NewCoin(types.GrantTokenDenom, math.OneInt()), sdk.NewCoin("zzz", math.OneInt())},
			err:       errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid fee token %s", "zzz"),
		},
		{
			name:      "invalid msg type",
			allowance: types.NewFeeAllowance(math.OneInt(), now),
			msg:       &blobtypes.MsgUpdateParams{},
			fee:       sdk.Coins{sdk.NewCoin(types.GrantTokenDenom, math.OneInt())},
			err:       errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid message type"),
		},
		{
			name:      "valid",
			allowance: types.NewFeeAllowance(math.OneInt(), now),
			fee:       sdk.Coins{sdk.NewCoin(types.GrantTokenDenom, math.OneInt())},
			msg:       &blobtypes.MsgPayForBlobs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remove, err := tt.allowance.Accept(context.TODO(), tt.fee, []sdk.Msg{tt.msg})
			if tt.err != nil {
				require.False(t, remove)
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.True(t, remove)
			require.NoError(t, err)
		})
	}
}

func TestFeeAllowanceValidateBasic(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		allowance types.FeeAllowance
		err       error
	}{
		{
			name:      "negative gas",
			allowance: types.NewFeeAllowance(math.NewInt(-1), now),
			err:       errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "gas must be positive: %d", -1),
		},
		{
			name:      "zero gas",
			allowance: types.NewFeeAllowance(math.ZeroInt(), now),
			err:       errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "gas must be positive: %d", 0),
		},
		{
			name:      "nil gas",
			allowance: types.NewFeeAllowance(math.Int{}, now),
			err:       errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "gas must be positive: nil"),
		},
		{
			name:      "valid gas",
			allowance: types.NewFeeAllowance(math.OneInt(), now),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.allowance.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
