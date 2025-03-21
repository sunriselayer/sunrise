package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) ClaimRewards(ctx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	totalRewards := sdk.NewCoins()

	for _, validatorAddrString := range msg.ValidatorAddresses {
		validatorAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(validatorAddrString)
		if err != nil {
			return nil, errorsmod.Wrap(err, "invalid validator address")
		}

		coins, err := k.Keeper.ClaimRewards(ctx, sender, validatorAddr)
		if err != nil {
			return nil, err
		}

		totalRewards = totalRewards.Add(coins...)
	}

	return &types.MsgClaimRewardsResponse{
		Amount: totalRewards,
	}, nil
}
