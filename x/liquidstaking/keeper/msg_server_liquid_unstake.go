package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stakingtypes "cosmossdk.io/x/staking/types"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k msgServer) LiquidUnstake(ctx context.Context, msg *types.MsgLiquidUnstake) (*types.MsgLiquidUnstakeResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// Claim rewards
	lstDenom := types.LiquidStakingTokenDenom(msg.ValidatorAddress)
	validatorAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	err = k.Keeper.ClaimRewards(ctx, sender, validatorAddr, lstDenom)
	if err != nil {
		return nil, err
	}

	// Get LST supply before burning
	lstSupplyOld := k.bankKeeper.GetSupply(ctx, lstDenom)

	// Send liquid staking token to module
	coins := sdk.NewCoins(sdk.NewCoin(lstDenom, msg.Amount))

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	// Burn liquid staking token
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	err = k.bankKeeper.BurnCoins(ctx, moduleAddr, coins)
	if err != nil {
		return nil, err
	}

	// Calculate unstake amount
	res, err := k.Environment.QueryRouterService.Invoke(ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: msg.Sender,
		ValidatorAddr: msg.ValidatorAddress,
	})
	if err != nil {
		return nil, err
	}
	queryDelegationResponse, ok := res.(*stakingtypes.QueryDelegationResponse)
	if !ok {
		return nil, sdkerrors.ErrInvalidRequest
	}
	stakedAmount := queryDelegationResponse.DelegationResponse.Balance.Amount
	outputAmount, err := types.CalculateLiquidUnstakeOutputAmount(stakedAmount, lstSupplyOld.Amount, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Undelegate
	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	output := sdk.NewCoin(params.BondDenom, outputAmount)

	res, err = k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgUndelegate{
		DelegatorAddress: msg.Sender,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           output,
	})
	if err != nil {
		return nil, err
	}
	undelegateResponse, ok := res.(*stakingtypes.MsgUndelegateResponse)
	if !ok {
		return nil, sdkerrors.ErrInvalidRequest
	}

	// Append Unstaking state
	_, err = k.AppendUnstaking(ctx, types.Unstaking{
		Address:        msg.Sender,
		CompletionTime: undelegateResponse.CompletionTime,
		Amount:         output,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgLiquidUnstakeResponse{
		CompletionTime: undelegateResponse.CompletionTime,
		Amount:         output,
	}, nil
}
