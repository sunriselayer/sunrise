package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ChallengeForFraud(ctx context.Context, msg *types.MsgChallengeForFraud) (*types.MsgChallengeForFraudResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	publishedData := k.GetPublishedData(ctx, msg.MetadataUri)
	if publishedData.Status != "vote_extension" {
		return nil, types.ErrCanNotOpenChallenge
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	if publishedData.Timestamp.Add(params.ChallengePeriod).Before(sdkCtx.BlockTime()) {
		return nil, types.ErrChallengePeriodIsOver
	}

	// Send collateral to module account
	if params.ChallengeCollateral.IsAllPositive() {
		sender := sdk.MustAccAddressFromBech32(msg.Sender)
		err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, publishedData.Collateral)
		if err != nil {
			return nil, err
		}
	}

	publishedData.Status = "challenge_for_fraud"
	publishedData.Challenger = msg.Sender
	publishedData.ChallengeTimestamp = sdkCtx.BlockTime()
	err = k.SetPublishedData(ctx, publishedData)
	if err != nil {
		return nil, err
	}

	err = sdkCtx.EventManager().EmitTypedEvent(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgChallengeForFraudResponse{}, nil
}
