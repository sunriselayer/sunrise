package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
	err = k.Keeper.ClaimRewards(ctx, msg.Sender, lstDenom)
	if err != nil {
		return nil, err
	}

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

	// TODO: Calculate unstake amount
	outputAmount := msg.Amount

	// Unstake
	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	res, err := k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgUndelegate{
		DelegatorAddress: msg.Sender,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           sdk.NewCoin(params.BondDenom, outputAmount),
	})
	if err != nil {
		return nil, err
	}
	response, ok := res.(*stakingtypes.MsgUndelegateResponse)
	if !ok {
		return nil, errorsmod.Wrap(errorsmod.ErrInvalidRequest, "invalid response")
	}

	// TODO: set Unstaking state

	return &types.MsgLiquidUnstakeResponse{}, nil
}
