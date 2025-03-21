package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	distributiontypes "cosmossdk.io/x/distribution/types"
	stakingtypes "cosmossdk.io/x/staking/types"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k Keeper) ValidateLastRewardHandlingTime(ctx context.Context, validatorAddr sdk.ValAddress) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	lastRewardHandlingTime, err := k.GetLastRewardHandlingTime(ctx, validatorAddr)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if sdkCtx.BlockTime().Before(lastRewardHandlingTime.Add(params.RewardPeriod)) {
		return nil
	}
	err = k.SetLastRewardHandlingTime(ctx, validatorAddr, sdkCtx.BlockTime())
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) HandleModuleAccountRewards(ctx context.Context) error {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	err := k.stakingKeeper.IterateDelegatorDelegations(ctx, moduleAddr, func(delegation stakingtypes.Delegation) (stop bool) {
		validatorAddr := delegation.ValidatorAddress

		err := k.HandleModuleAccountRewardsByValidator(ctx, validatorAddr)
		if err != nil {
			k.Logger.Error("failed to handle module account rewards by validator", "validator", validatorAddr, "error", err)
		}

		return false
	})
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) HandleModuleAccountRewardsByValidator(ctx context.Context, validatorAddr string) error {
	validatorAddrBytes, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(validatorAddr)
	if err != nil {
		return err
	}

	// Validate last reward handling time to mitigate the load
	err = k.ValidateLastRewardHandlingTime(ctx, validatorAddrBytes)
	if err != nil {
		return err
	}

	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	rewardSaverAddr := types.RewardSaverAddress(validatorAddr)

	// Withdraw delegator reward as module address
	res, err := k.Environment.MsgRouterService.Invoke(ctx, &distributiontypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: validatorAddr,
	})
	if err != nil {
		return err
	}

	response, ok := res.(*distributiontypes.MsgWithdrawDelegatorRewardResponse)
	if !ok {
		return sdkerrors.ErrInvalidRequest
	}

	if response.Amount.IsZero() {
		return nil
	}

	// Transfer to reward saver address
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, rewardSaverAddr, response.Amount)
	if err != nil {
		return err
	}

	// Get total share
	totalShare := k.GetTotalShare(ctx, validatorAddr)

	if totalShare.IsZero() {
		return nil
	}

	// Iterate all rewards
	// Multiplier_new = Multiplier_old + (Reward_new) / totalShare
	for _, coin := range response.Amount {
		multiplierOld, err := k.GetRewardMultiplier(ctx, validatorAddrBytes, coin.Denom)
		if err != nil {
			return err
		}

		multiplierNew, err := types.CalculateRewardMultiplierNew(multiplierOld, coin.Amount, totalShare)
		if err != nil {
			return err
		}

		err = k.SetRewardMultiplier(ctx, validatorAddrBytes, coin.Denom, multiplierNew)
		if err != nil {
			return err
		}
	}

	return nil
}
