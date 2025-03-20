package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k Keeper) GetRewardMultiplier(ctx context.Context, validatorAddr sdk.ValAddress, denom string) (math.Dec, error) {
	rewardMultiplierString, err := k.RewardMultiplier.Get(ctx, collections.Join([]byte(validatorAddr), denom))
	if err == collections.ErrNotFound {
		return math.NewDecFromInt64(0), nil
	}

	if err != nil {
		return math.Dec{}, err
	}
	return math.NewDecFromString(rewardMultiplierString)
}

func (k Keeper) SetRewardMultiplier(ctx context.Context, validatorAddr sdk.ValAddress, denom string, value math.Dec) error {
	err := k.RewardMultiplier.Set(ctx, collections.Join([]byte(validatorAddr), denom), value.String())
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetUserLastRewardMultiplier(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress, denom string) (math.Dec, error) {
	userLastRewardMultiplierString, err := k.UsersLastRewardMultiplier.Get(ctx, collections.Join3(sender, []byte(validatorAddr), denom))
	if err == collections.ErrNotFound {
		return math.NewDecFromInt64(0), nil
	}

	if err != nil {
		return math.Dec{}, err
	}
	return math.NewDecFromString(userLastRewardMultiplierString)
}

func (k Keeper) SetUserLastRewardMultiplier(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress, denom string, value math.Dec) error {
	err := k.UsersLastRewardMultiplier.Set(ctx, collections.Join3(sender, []byte(validatorAddr), denom), value.String())
	if err != nil {
		return err
	}
	return nil
}

// ClaimRewards claims rewards from a validator
// reward = (rewardMultiplier - userLastRewardMultiplier) * lstCoin.Amount
func (k Keeper) ClaimRewards(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress, denom string) error {
	// LST info
	lstDenom := types.LiquidStakingTokenDenom(validatorAddr.String())
	lstCoin := k.bankKeeper.GetBalance(ctx, sender, lstDenom)

	rewardMultiplier, err := k.GetRewardMultiplier(ctx, validatorAddr, denom)
	if err != nil {
		return err
	}

	userLastRewardMultiplier, err := k.GetUserLastRewardMultiplier(ctx, sender, validatorAddr, denom)
	if err != nil {
		return err
	}

	sub, err := rewardMultiplier.Sub(userLastRewardMultiplier)
	if err != nil {
		return err
	}

	rewardAmountDec, err := sub.Mul(math.NewDecFromInt64(lstCoin.Amount.Int64()))
	if err != nil {
		return err
	}
	rewardAmount, err := rewardAmountDec.SdkIntTrim()
	if err != nil {
		return err
	}

	rewardCoin := sdk.NewCoins(sdk.NewCoin(denom, rewardAmount))

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.RewardSaverModuleAccount(validatorAddr.String()), sender, rewardCoin)
	if err != nil {
		return err
	}

	return nil
}
