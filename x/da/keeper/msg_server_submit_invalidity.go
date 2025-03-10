package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SubmitInvalidity(ctx context.Context, msg *types.MsgSubmitInvalidity) (*types.MsgSubmitInvalidityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}
	// check number of indices
	if len(msg.Indices) == 0 {
		return nil, types.ErrInvalidIndices
	}
	// end static validation

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	publishedData, found := k.GetPublishedData(ctx, msg.MetadataUri)
	if !found {
		return nil, types.ErrDataNotFound
	}
	if publishedData.Status != types.Status_STATUS_CHALLENGE_PERIOD {
		return nil, types.ErrNotInChallengePeriod
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	if publishedData.Timestamp.Add(params.ChallengePeriod).Before(sdkCtx.BlockTime()) {
		return nil, types.ErrChallengePeriodIsOver
	}

	// Send collateral to module account
	if publishedData.SubmitInvalidityCollateral.IsAllPositive() {
		sender := sdk.MustAccAddressFromBech32(msg.Sender)
		err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, publishedData.SubmitInvalidityCollateral)
		if err != nil {
			return nil, err
		}
	}

	// save invalidity in the storage
	err = k.SetInvalidity(ctx, types.Invalidity{
		MetadataUri: msg.MetadataUri,
		Sender:      msg.Sender,
		Indices:     msg.Indices,
	})
	if err != nil {
		return nil, err
	}

	err = sdkCtx.EventManager().EmitTypedEvent(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitInvalidityResponse{}, nil
}
