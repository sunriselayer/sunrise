package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k msgServer) InitLockupAccount(ctx context.Context, msg *types.MsgInitLockupAccount) (*types.MsgInitLockupAccountResponse, error) {
	err := k.InitLockupAccountFromMsg(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgInitLockupAccountResponse{}, nil
}
