package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k msgServer) InitLockupAccount(ctx context.Context, msg *types.MsgInitLockupAccount) (*types.MsgInitLockupAccountResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	// TODO: Handle the message

	return &types.MsgInitLockupAccountResponse{}, nil
}
