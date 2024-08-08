package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) ChallengeForFraud(goCtx context.Context, msg *types.MsgChallengeForFraud) (*types.MsgChallengeForFraudResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	publishedData := k.GetPublishedData(ctx, msg.MetadataUri)
	if publishedData.Status != "vote_extension" {
		return nil, types.ErrCanNotOpenChallenge
	}

	params := k.GetParams(ctx)
	if publishedData.Timestamp.Add(params.ChallengePeriod).Before(ctx.BlockTime()) {
		return nil, types.ErrChallengePeriodIsOver
	}

	// Send collateral to module account
	if params.ChallengeCollateral.IsAllPositive() {
		sender := sdk.MustAccAddressFromBech32(msg.Sender)
		err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, params.ChallengeCollateral)
		if err != nil {
			return nil, err
		}
	}

	publishedData.Status = "challenge_for_fraud"
	publishedData.Challenger = msg.Sender
	publishedData.Collateral = params.ChallengeCollateral
	publishedData.ChallengeTimestamp = ctx.BlockTime()
	err := k.SetPublishedData(ctx, publishedData)
	if err != nil {
		return nil, err
	}

	err = ctx.EventManager().EmitTypedEvent(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgChallengeForFraudResponse{}, nil
}
