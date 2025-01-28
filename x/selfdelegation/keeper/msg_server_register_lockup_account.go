package keeper

import (
	"context"

	"cosmossdk.io/errors"

	"github.com/sunriselayer/sunrise/x/selfdelegation/types"
)

func (k msgServer) RegisterLockupAccount(ctx context.Context, msg *types.MsgRegisterLockupAccount) (*types.MsgRegisterLockupAccountResponse, error) {
	lockupBytes, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errors.Wrap(err, "invalid lockup account address")
	}
	ownerBytes, err := k.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "invalid owner address")
	}
	err = k.LockupAccounts.Set(ctx, lockupBytes, ownerBytes)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterLockupAccountResponse{}, nil
}
