package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "cosmossdk.io/x/staking/types"
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

	// Prepare fee and bond coin
	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	feeCoin := sdk.NewCoin(params.FeeDenom, msg.Amount)
	bondCoin := sdk.NewCoin(params.BondDenom, msg.Amount)

	// Send fee coin to module

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return nil, err
	}

	// Convert fee denom to bond denom
	err = k.tokenConverterKeeper.ConvertReverse(ctx, msg.Amount, sender)
	if err != nil {
		return nil, err
	}

	//  Stake
	_, err = k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: msg.Sender,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           bondCoin,
	})
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
