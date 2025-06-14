package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func (k msgServer) ClaimYield(ctx context.Context, msg *types.MsgClaimYield) (*types.MsgClaimYieldResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgClaimYieldResponse{}, nil
}
