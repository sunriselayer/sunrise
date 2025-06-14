package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func (k msgServer) AddYields(ctx context.Context, msg *types.MsgAddYields) (*types.MsgAddYieldsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Admin); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgAddYieldsResponse{}, nil
}
