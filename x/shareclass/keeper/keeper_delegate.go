package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "cosmossdk.io/x/staking/types"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k Keeper) ConvertAndDelegate(ctx context.Context, sender sdk.AccAddress, validatorAddr string, amount math.Int) error {
	// Prepare fee and bond coin
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return err
	}
	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return err
	}
	bondCoin := sdk.NewCoin(bondDenom, amount)
	feeCoin := sdk.NewCoin(feeDenom, amount)

	// Send fee coin to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return err
	}

	// Convert fee denom to bond denom
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	err = k.tokenConverterKeeper.ConvertReverse(ctx, amount, moduleAddr)
	if err != nil {
		return err
	}

	// Stake
	_, err = k.StakingMsgServer.Delegate(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: moduleAddr.String(),
		ValidatorAddress: validatorAddr,
		Amount:           bondCoin,
	})
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetTotalStakedAmount(ctx context.Context, validatorAddr string) (math.Int, error) {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	res, err := k.StakingQueryServer.Delegation(ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: moduleAddr.String(),
		ValidatorAddr: validatorAddr,
	})
	if err != nil {
		return math.Int{}, err
	}

	stakedAmount := res.DelegationResponse.Balance.Amount

	return stakedAmount, nil
}
