package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func (k msgServer) UpdateAdmin(ctx context.Context, msg *types.MsgUpdateAdmin) (*types.MsgUpdateAdminResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Admin); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgUpdateAdminResponse{}, nil
}
