package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k msgServer) LiquidUnstake(ctx context.Context, msg *types.MsgLiquidUnstake) (*types.MsgLiquidUnstakeResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// Claim rewards
	err = k.Keeper.AutoCompoundModuleAccountRewards(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	// Send liquid staking token to module
	denom := types.LiquidStakingTokenDenom(msg.ValidatorAddress)
	coins := sdk.NewCoins(sdk.NewCoin(denom, msg.Amount))

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

	// TODO: Unstake
	outputAmount := msg.Amount

	// Convert bond denom to fee denom
	err = k.tokenConverterKeeper.Convert(ctx, outputAmount, sender)
	if err != nil {
		return nil, err
	}

	// Send fee coin to sender
	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	feeCoin := sdk.NewCoin(params.FeeDenom, outputAmount)

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(feeCoin))
	if err != nil {
		return nil, err
	}

	return &types.MsgLiquidUnstakeResponse{}, nil
}
