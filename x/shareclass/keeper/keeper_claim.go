package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k Keeper) GetRewardMultiplier(ctx context.Context, validatorAddr sdk.ValAddress, denom string) (math.LegacyDec, error) {
	rewardMultiplierString, err := k.RewardMultiplier.Get(ctx, collections.Join([]byte(validatorAddr), denom))
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return math.LegacyZeroDec(), nil
		}
		return math.LegacyDec{}, err
	}
	return math.LegacyNewDecFromStr(rewardMultiplierString)
}

func (k Keeper) SetRewardMultiplier(ctx context.Context, validatorAddr sdk.ValAddress, denom string, value math.LegacyDec) error {
	err := k.RewardMultiplier.Set(ctx, collections.Join([]byte(validatorAddr), denom), value.String())
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetUserLastRewardMultiplier(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress, denom string) (math.LegacyDec, error) {
	userLastRewardMultiplierString, err := k.UsersLastRewardMultiplier.Get(ctx, collections.Join3(sender, []byte(validatorAddr), denom))
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return math.LegacyZeroDec(), nil
		}
		return math.LegacyDec{}, err
	}
	return math.LegacyNewDecFromStr(userLastRewardMultiplierString)
}

func (k Keeper) SetUserLastRewardMultiplier(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress, denom string, value math.LegacyDec) error {
	err := k.UsersLastRewardMultiplier.Set(ctx, collections.Join3(sender, []byte(validatorAddr), denom), value.String())
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) ClaimRewards(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress) (sdk.Coins, error) {
	total, err := k.GetClaimableRewards(ctx, sender, validatorAddr)
	if err != nil {
		return nil, err
	}

	if total.IsZero() {
		return total, nil
	}

	err = k.bankKeeper.SendCoins(ctx, types.RewardSaverAddress(validatorAddr.String()), sender, total)
	if err != nil {
		return nil, err
	}

	// Update user's last reward multiplier
	for _, coin := range total {
		rewardMultiplier, err := k.GetRewardMultiplier(ctx, validatorAddr, coin.Denom)
		if err != nil {
			return nil, err
		}
		err = k.SetUserLastRewardMultiplier(ctx, sender, validatorAddr, coin.Denom, rewardMultiplier)
		if err != nil {
			return nil, err
		}
	}

	return total, nil
}

func (k Keeper) GetClaimableRewards(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress) (sdk.Coins, error) {
	coins := k.bankKeeper.GetAllBalances(ctx, types.RewardSaverAddress(validatorAddr.String()))
	total := sdk.NewCoins()

	for _, coin := range coins {
		amount, err := k.GetClaimableRewardsByDenom(ctx, sender, validatorAddr, coin.Denom)
		if err != nil {
			return nil, err
		}
		total = total.Add(sdk.NewCoin(coin.Denom, amount))
	}

	return total, nil
}

// GetClaimableRewardsByDenom claims rewards from a validator
// reward = (rewardMultiplier - userLastRewardMultiplier) * share.Amount
func (k Keeper) GetClaimableRewardsByDenom(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress, denom string) (math.Int, error) {
	// Get the share
	share := k.GetShare(ctx, sender, validatorAddr.String())

	// Get the reward multiplier
	rewardMultiplier, err := k.GetRewardMultiplier(ctx, validatorAddr, denom)
	if err != nil {
		return math.Int{}, err
	}

	userLastRewardMultiplier, err := k.GetUserLastRewardMultiplier(ctx, sender, validatorAddr, denom)
	if err != nil {
		return math.Int{}, err
	}

	rewardAmount, err := types.CalculateReward(rewardMultiplier, userLastRewardMultiplier, share)
	if err != nil {
		return math.Int{}, err
	}

	return rewardAmount, nil
}
