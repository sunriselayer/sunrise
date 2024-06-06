package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) CollectFees(goCtx context.Context, msg *types.MsgCollectFees) (*types.MsgCollectFeesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	totalCollectedFees := sdk.NewCoins()
	for _, positionId := range msg.PositionIds {
		collectedFees, err := k.Keeper.collectFees(ctx, sender, positionId)
		if err != nil {
			return nil, err
		}
		totalCollectedFees = totalCollectedFees.Add(collectedFees...)
	}

	return &types.MsgCollectFeesResponse{
		CollectedFees: totalCollectedFees,
	}, nil
}
