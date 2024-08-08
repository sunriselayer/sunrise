package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) PublishData(goCtx context.Context, msg *types.MsgPublishData) (*types.MsgPublishDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.SetPublishedData(ctx, types.PublishedData{
		MetadataUri:       msg.MetadataUri,
		ShardDoubleHashes: msg.ShardDoubleHashes,
		Timestamp:         ctx.BlockTime(),
		Status:            "msg_server",
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgPublishDataResponse{}, nil
}
