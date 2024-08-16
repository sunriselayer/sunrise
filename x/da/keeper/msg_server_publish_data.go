package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) PublishData(goCtx context.Context, msg *types.MsgPublishData) (*types.MsgPublishDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	err := k.SetPublishedData(ctx, types.PublishedData{
		Publisher:         msg.Sender,
		MetadataUri:       msg.MetadataUri,
		ParityShardCount:  msg.ParityShardCount,
		ShardDoubleHashes: msg.ShardDoubleHashes,
		Collateral:        params.ChallengeCollateral,
		Timestamp:         ctx.BlockTime(),
		Status:            "msg_server",
	})
	if err != nil {
		return nil, err
	}

	// Send collateral to module account
	if params.ChallengeCollateral.IsAllPositive() {
		sender := sdk.MustAccAddressFromBech32(msg.Sender)
		err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, params.ChallengeCollateral)
		if err != nil {
			return nil, err
		}
	}

	err = ctx.EventManager().EmitTypedEvent(msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgPublishDataResponse{}, nil
}
