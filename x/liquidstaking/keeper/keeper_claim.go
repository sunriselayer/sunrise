package keeper

import (
	"context"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	distributiontypes "cosmossdk.io/x/distribution/types"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k Keeper) AutoCompoundModuleAccountRewards(ctx context.Context, validatorAddr string) error {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	rewardSaverModuleAddr := k.accountKeeper.GetModuleAddress(types.RewardSaverModuleAccount())

	res, err := k.Environment.MsgRouterService.Invoke(ctx, &distributiontypes.MsgSetWithdrawAddress{
		DelegatorAddress: moduleAddr.String(),
		WithdrawAddress:  rewardSaverModuleAddr.String(),
	})
	if err != nil {
		return err
	}

	res, err = k.Environment.MsgRouterService.Invoke(ctx, &distributiontypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: validatorAddr,
	})
	if err != nil {
		return err
	}

	response, ok := res.(*distributiontypes.MsgWithdrawDelegatorRewardResponse)
	if !ok {
		return errorsmod.Wrap(errorsmod.ErrInvalidRequest, "invalid response")
	}

	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return err
	}

	ok, feeCoin := response.Amount.Find(params.FeeDenom)
	if !ok {
		return errorsmod.Wrap(errorsmod.ErrInvalidRequest, "invalid response")
	}

	if feeCoin.Amount.IsZero() {
		return nil
	}

	res, err = k.Environment.MsgRouterService.Invoke(ctx, &types.MsgLiquidStake{
		Sender:           rewardSaverModuleAddr.String(),
		ValidatorAddress: validatorAddr,
		Amount:           feeCoin.Amount,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) ClaimRewards(ctx context.Context, sender string, denom string) error {

}
