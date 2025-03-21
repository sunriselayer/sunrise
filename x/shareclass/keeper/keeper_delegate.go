package keeper

import (
	"context"

	"cosmossdk.io/math"
	stakingtypes "cosmossdk.io/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k Keeper) ConvertAndDelegate(ctx context.Context, sender sdk.AccAddress, validatorAddr string, amount math.Int) error {
	// Prepare fee and bond coin
	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return err
	}
	feeCoin := sdk.NewCoin(params.FeeDenom, amount)
	bondCoin := sdk.NewCoin(params.BondDenom, amount)

	// Send fee coin to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return err
	}

	// Convert fee denom to bond denom
	err = k.tokenConverterKeeper.ConvertReverse(ctx, amount, sender)
	if err != nil {
		return err
	}

	// Stake
	_, err = k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sender.String(),
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

	res, err := k.Environment.QueryRouterService.Invoke(ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: moduleAddr.String(),
		ValidatorAddr: validatorAddr,
	})
	if err != nil {
		return math.Int{}, err
	}
	queryDelegationResponse, ok := res.(*stakingtypes.QueryDelegationResponse)
	if !ok {
		return math.Int{}, sdkerrors.ErrInvalidRequest
	}
	stakedAmount := queryDelegationResponse.DelegationResponse.Balance.Amount

	return stakedAmount, nil
}
