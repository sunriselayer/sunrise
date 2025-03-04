package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ClaimRewards(ctx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	if len(msg.PositionIds) == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "position ids cannot be empty")
	}
	// end static validation

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	totalCollectedFees := sdk.NewCoins()
	for _, positionId := range msg.PositionIds {
		collectedFees, err := k.Keeper.collectFees(sdkCtx, sender, positionId)
		if err != nil {
			return nil, err
		}
		totalCollectedFees = totalCollectedFees.Add(collectedFees...)
	}

	return &types.MsgClaimRewardsResponse{
		CollectedFees: totalCollectedFees,
	}, nil
}
