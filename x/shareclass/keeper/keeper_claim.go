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
	totalClaimable, rewardMultipliers, err := k.GetClaimableRewards(ctx, sender, validatorAddr)
	if err != nil {
		return nil, err
	}

	// Send the total claimable rewards if any.
	if !totalClaimable.IsZero() {
		err := k.bankKeeper.SendCoins(ctx, types.RewardSaverAddress(validatorAddr.String()), sender, totalClaimable)
		if err != nil {
			return nil, err
		}
	}

	// Update the user's last reward multiplier for all denoms using the multipliers fetched during calculation.
	for denom, multiplier := range rewardMultipliers {
		err := k.SetUserLastRewardMultiplier(ctx, sender, validatorAddr, denom, multiplier)
		if err != nil {
			return nil, err
		}
	}

	return totalClaimable, nil
}

// GetClaimableRewards calculates the total claimable rewards and returns them along with the reward multipliers used.
func (k Keeper) GetClaimableRewards(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress) (sdk.Coins, map[string]math.LegacyDec, error) {
	// Get all denoms from the reward saver address.
	allRewardCoins := k.bankKeeper.GetAllBalances(ctx, types.RewardSaverAddress(validatorAddr.String()))

	// Get the user's share.
	share := k.GetShare(ctx, sender, validatorAddr.String())

	// Prepare to calculate the total claimable amount and store the latest reward multipliers.
	totalClaimable := sdk.NewCoins()
	rewardMultipliers := make(map[string]math.LegacyDec)

	// Calculate claimable rewards for each denom and store the multiplier.
	for _, coin := range allRewardCoins {
		denom := coin.Denom

		rewardAmount, rewardMultiplier, err := k.GetClaimableRewardsByDenom(ctx, sender, validatorAddr, denom, share)
		if err != nil {
			return nil, nil, err
		}

		rewardMultipliers[denom] = rewardMultiplier

		// Add to total if there's a positive amount to claim.
		if rewardAmount.IsPositive() {
			totalClaimable = totalClaimable.Add(sdk.NewCoin(denom, rewardAmount))
		}
	}

	return totalClaimable, rewardMultipliers, nil
}

// GetClaimableRewardsByDenom calculates the reward for a single denomination.
func (k Keeper) GetClaimableRewardsByDenom(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress, denom string, share math.Int) (math.Int, math.LegacyDec, error) {
	// Fetch reward multiplier once.
	rewardMultiplier, err := k.GetRewardMultiplier(ctx, validatorAddr, denom)
	if err != nil {
		return math.Int{}, math.LegacyDec{}, err
	}

	// Get the user's last known multiplier.
	userLastRewardMultiplier, err := k.GetUserLastRewardMultiplier(ctx, sender, validatorAddr, denom)
	if err != nil {
		return math.Int{}, math.LegacyDec{}, err
	}

	// Calculate the reward amount.
	rewardAmount, err := types.CalculateReward(rewardMultiplier, userLastRewardMultiplier, share)
	if err != nil {
		return math.Int{}, math.LegacyDec{}, err
	}

	return rewardAmount, rewardMultiplier, nil
}
