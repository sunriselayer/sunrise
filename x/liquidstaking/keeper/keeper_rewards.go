package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	distributiontypes "cosmossdk.io/x/distribution/types"
	stakingtypes "cosmossdk.io/x/staking/types"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
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

func (k Keeper) GetRewardSaverAddress(ctx context.Context, validatorAddr string) sdk.AccAddress {
	rewardSaverAccount := types.RewardSaverModuleAccount(validatorAddr)

	return k.accountKeeper.GetModuleAddress(rewardSaverAccount)
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
	rewardSaverAddr := k.GetRewardSaverAddress(ctx, validatorAddr)

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

	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return err
	}

	// Get LST info
	lstDenom := types.LiquidStakingTokenDenom(validatorAddr)
	lstSupplyOld := k.bankKeeper.GetSupply(ctx, lstDenom)

	// Convert fee coin to LST
	ok, feeCoin := response.Amount.Find(params.FeeDenom)
	if ok {
		stakedAmount, err := k.GetStakedAmount(ctx, validatorAddr)
		if err != nil {
			return err
		}
		compensation, distribution := types.CalculateSlashingCompensation(stakedAmount, lstSupplyOld.Amount, feeCoin.Amount)

		// For slashing compensation
		if compensation.IsPositive() {
			err = k.ConvertAndDelegate(ctx, rewardSaverAddr, validatorAddr, compensation)
			if err != nil {
				return err
			}
		}

		// For LST distribution
		if distribution.IsPositive() {
			_, err = k.Environment.MsgRouterService.Invoke(ctx, &types.MsgLiquidStake{
				Sender:           rewardSaverAddr.String(),
				ValidatorAddress: validatorAddr,
				Amount:           distribution,
			})
			if err != nil {
				return err
			}
		}
	}

	// Get LST info
	lstSupplyNew := k.bankKeeper.GetSupply(ctx, lstDenom)

	if lstSupplyNew.IsZero() {
		return nil
	}

	// Iterate all rewards
	// Multiplier_new = Multiplier_old + (Reward_new) / Supply_LST_new
	// Supply_LST_new = Supply_LST_old + Reward_LST_new
	for _, coin := range response.Amount {
		multiplierDenom := coin.Denom
		if multiplierDenom == params.FeeDenom {
			multiplierDenom = lstDenom
		}

		multiplierOld, err := k.GetRewardMultiplier(ctx, validatorAddrBytes, multiplierDenom)
		if err != nil {
			return err
		}

		multiplierNew, err := types.CalculateRewardMultiplierNew(multiplierOld, coin.Amount, lstSupplyNew.Amount)
		if err != nil {
			return err
		}

		err = k.SetRewardMultiplier(ctx, validatorAddrBytes, multiplierDenom, multiplierNew)
		if err != nil {
			return err
		}
	}

	return nil
}
