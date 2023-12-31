package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"sunrise/x/blobstream/types"
)

func (k msgServer) RegisterEvmAddress(goCtx context.Context, msg *types.MsgRegisterEvmAddress) (*types.MsgRegisterEvmAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgRegisterEvmAddressResponse{}, nil
}
