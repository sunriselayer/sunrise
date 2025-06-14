package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func (k msgServer) ClaimYields(ctx context.Context, msg *types.MsgClaimYields) (*types.MsgClaimYieldsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgClaimYieldsResponse{}, nil
}
