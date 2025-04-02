package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) ClaimRewards(ctx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	validatorAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(msg.Validator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	rewards, err := k.Keeper.ClaimRewards(ctx, sender, validatorAddr)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimRewardsResponse{
		Amount: rewards,
	}, nil
}
