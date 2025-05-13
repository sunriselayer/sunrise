package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k Keeper) StartNewEpoch(ctx context.Context, msg *types.MsgStartNewEpoch) (*types.MsgStartNewEpochResponse, error) {
	_, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Create a new `Epoch` if the last `Epoch` has ended or the first `Epoch` has not been created.
	lastEpoch, found, err := k.GetLastEpoch(sdkCtx)
	if err != nil {
		return nil, err
	}
	var epochId uint64
	if !found {
		epochId = 1
	} else {
		epochId = lastEpoch.Id + 1

		if sdkCtx.BlockHeight() < lastEpoch.EndBlock {
			return nil, types.ErrEpochNotEnded
		}
	}

	err = k.CreateEpoch(sdkCtx, epochId)
	if err != nil {
		return nil, err
	}

	return &types.MsgStartNewEpochResponse{}, nil
}
