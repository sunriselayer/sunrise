package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	distributiontypes "cosmossdk.io/x/distribution/types"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k Keeper) HandleModuleAccountRewards(ctx context.Context) error {
	// TODO: iterate all validators which the module account delegates to

	return nil
}

func (k Keeper) GetRewardSaverAddress(ctx context.Context, validatorAddr string) sdk.AccAddress {
	rewardSaverAccount := types.RewardSaverModuleAccount(validatorAddr)

	return k.accountKeeper.GetModuleAddress(rewardSaverAccount)
}

func (k Keeper) HandleModuleAccountRewardsByValidator(ctx context.Context, validatorAddr sdk.ValAddress) error {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	validatorAddrStr := validatorAddr.String()
	rewardSaverAddr := k.GetRewardSaverAddress(ctx, validatorAddrStr)

	// Withdraw delegator reward as module address
	res, err := k.Environment.MsgRouterService.Invoke(ctx, &distributiontypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: validatorAddrStr,
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
	lstDenom := types.LiquidStakingTokenDenom(validatorAddrStr)
	lstSupplyOld := k.bankKeeper.GetSupply(ctx, lstDenom)

	// Convert fee coin to LST
	ok, feeCoin := response.Amount.Find(params.FeeDenom)
	if ok {
		stakedAmount, err := k.GetStakedAmount(ctx, validatorAddrStr)
		if err != nil {
			return err
		}
		compensation, distribution := types.CalculateSlashingCompensation(stakedAmount, lstSupplyOld.Amount, feeCoin.Amount)

		if compensation.IsPositive() {
			err = k.ConvertAndDelegate(ctx, rewardSaverAddr, validatorAddrStr, compensation)
			if err != nil {
				return err
			}
		}
		if distribution.IsPositive() {
			_, err = k.Environment.MsgRouterService.Invoke(ctx, &types.MsgLiquidStake{
				Sender:           rewardSaverAddr.String(),
				ValidatorAddress: validatorAddrStr,
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

		multiplierOld, err := k.GetRewardMultiplier(ctx, validatorAddr, multiplierDenom)
		if err != nil {
			return err
		}

		multiplierNew, err := types.CalculateRewardMultiplierNew(multiplierOld, coin.Amount, lstSupplyNew.Amount)
		if err != nil {
			return err
		}

		err = k.SetRewardMultiplier(ctx, validatorAddr, multiplierDenom, multiplierNew)
		if err != nil {
			return err
		}
	}

	return nil
}
