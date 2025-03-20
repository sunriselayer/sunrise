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
	lstDenom := types.LiquidStakingTokenDenom(msg.ValidatorAddress)
	validatorAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	err = k.Keeper.ClaimRewards(ctx, sender, validatorAddr, lstDenom)
	if err != nil {
		return nil, err
	}

	// Convert and delegate
	err = k.ConvertAndDelegate(ctx, sender, msg.ValidatorAddress, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Mint liquid staking token
	coins := sdk.NewCoins(sdk.NewCoin(lstDenom, msg.Amount))

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
