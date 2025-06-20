package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func (k msgServer) Borrow(ctx context.Context, msg *types.MsgBorrow) (*types.MsgBorrowResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgBorrowResponse{}, nil
}
