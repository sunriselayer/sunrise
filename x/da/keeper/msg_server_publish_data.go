package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) PublishData(ctx context.Context, msg *types.MsgPublishData) (*types.MsgPublishDataResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgPublishDataResponse{}, nil
}
