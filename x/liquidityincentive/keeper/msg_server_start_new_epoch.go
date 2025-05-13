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
	var lastEpochId uint64
	var newEpochId uint64
	if !found {
		lastEpochId = 0
		newEpochId = 1
	} else {
		lastEpochId = lastEpoch.Id
		newEpochId = lastEpochId + 1

		if sdkCtx.BlockHeight() < lastEpoch.EndBlock {
			return nil, types.ErrEpochNotEnded
		}
	}

	err = k.CreateEpoch(sdkCtx, lastEpochId, newEpochId)
	if err != nil {
		return nil, err
	}

	// Finalize bribe for new epoch & remove old epochs
	if err := k.FinalizeBribeForEpoch(sdkCtx, newEpochId); err != nil {
		return nil, err
	}

	return &types.MsgStartNewEpochResponse{}, nil
}
