package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) Convert(ctx context.Context, msg *types.MsgConvert) (*types.MsgConvertResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	if !msg.Amount.IsPositive() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be positive")
	}

	// end static validation
	address, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.Convert(ctx, msg.Amount, address); err != nil {
		return nil, err
	}

	return &types.MsgConvertResponse{}, nil
}
