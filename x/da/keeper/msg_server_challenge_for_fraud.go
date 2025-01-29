package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ChallengeForFraud(ctx context.Context, msg *types.MsgChallengeForFraud) (*types.MsgChallengeForFraudResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}
	// end static validation

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	publishedData := k.GetPublishedData(ctx, msg.MetadataUri)
	if publishedData.Status != types.Status_STATUS_CHALLENGE_PERIOD {
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
	if publishedData.Collateral.IsAllPositive() {
		sender := sdk.MustAccAddressFromBech32(msg.Sender)
		err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, publishedData.Collateral)
		if err != nil {
			return nil, err
		}
	}

	publishedData.Status = types.Status_STATUS_CHALLENGING
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
