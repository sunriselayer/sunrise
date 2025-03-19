package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k msgServer) LiquidStake(ctx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// Claim rewards
	err = k.Keeper.AutoCompoundModuleAccountRewards(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	// Send fee coin to module
	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	feeCoin := sdk.NewCoin(params.FeeDenom, msg.Amount)

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return nil, err
	}

	// Convert fee denom to bond denom
	err = k.tokenConverterKeeper.ConvertReverse(ctx, msg.Amount, sender)
	if err != nil {
		return nil, err
	}

	// TODO: Stake

	// Mint liquid staking token
	denom := types.LiquidStakingTokenDenom(msg.ValidatorAddress)
	coins := sdk.NewCoins(sdk.NewCoin(denom, msg.Amount))

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	// Send liquid staking token to sender
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coins)
	if err != nil {
		return nil, err
	}

	return &types.MsgLiquidStakeResponse{}, nil
}
