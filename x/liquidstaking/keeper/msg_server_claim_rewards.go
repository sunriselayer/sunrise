package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k msgServer) ClaimRewards(ctx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	_, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	rewardSaverModuleAddr := k.accountKeeper.GetModuleAddress(types.RewardSaverModuleAccount())
	coins := k.bankKeeper.SpendableCoins(ctx, rewardSaverModuleAddr)

	for _, coin := range coins {
		err = k.Keeper.ClaimRewards(ctx, msg.Sender, coin.Denom)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgClaimRewardsResponse{}, nil
}
