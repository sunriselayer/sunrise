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
	if !found {
		err := k.CreateEpoch(sdkCtx, 0, 1)
		if err != nil {
			return nil, err
		}

		return &types.MsgStartNewEpochResponse{}, nil
	}

	if sdkCtx.BlockHeight() < lastEpoch.EndBlock {
		return nil, types.ErrEpochNotEnded
	}
	// End current epoch and start new one
	if err := k.FinalizeBribeForEpoch(sdkCtx); err != nil {
		return nil, err
	}

	err = k.CreateEpoch(sdkCtx, lastEpoch.Id, lastEpoch.Id+1)
	if err != nil {
		return nil, err
	}
	// remove old epoch and gauges
	epochs, err := k.GetAllEpoch(sdkCtx)
	if err != nil {
		return nil, err
	}
	params, err := k.Params.Get(sdkCtx)
	if err != nil {
		return nil, err
	}
	for len(epochs) > int(params.BribeClaimEpochs)+1 {
		epochToRemove := epochs[0]
		err := k.RemoveEpoch(sdkCtx, epochToRemove.Id)
		if err != nil {
			return nil, err
		}
		for _, gauge := range epochToRemove.Gauges {
			err := k.RemoveGauge(sdkCtx, gauge.PreviousEpochId, gauge.PoolId)
			if err != nil {
				return nil, err
			}
		}
		// Remove the processed epoch from the slice to correctly check the condition in the next iteration
		epochs = epochs[1:]
	}

	// Event is emitted in CreateEpoch

	return &types.MsgStartNewEpochResponse{}, nil
}
