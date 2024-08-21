package keeper

import (
	"context"

	// sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/blob/types"
)

func (k msgServer) PayForBlobs(goCtx context.Context, msg *types.MsgPayForBlobs) (*types.MsgPayForBlobsResponse, error) {
	// 	ctx := sdk.UnwrapSDKContext(goCtx)

	// 	gasToConsume := types.GasToConsume(msg.BlobSizes, k.GasPerBlobByte(ctx))
	// 	ctx.GasMeter().ConsumeGas(gasToConsume, payForBlobGasDescriptor)

	// 	err := ctx.EventManager().EmitTypedEvent(
	// 		types.NewPayForBlobsEvent(msg.Signer, msg.BlobSizes, msg.Namespaces),
	// 	)
	// 	if err != nil {
	// 		return &types.MsgPayForBlobsResponse{}, err
	// 	}

	return &types.MsgPayForBlobsResponse{}, nil
}
