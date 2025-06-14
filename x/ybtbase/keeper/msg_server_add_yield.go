package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func (k msgServer) AddYield(ctx context.Context, msg *types.MsgAddYield) (*types.MsgAddYieldResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Admin); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgAddYieldResponse{}, nil
}
