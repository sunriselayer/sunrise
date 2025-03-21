package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) NonVotingDelegate(ctx context.Context, msg *types.MsgNonVotingDelegate) (*types.MsgNonVotingDelegateResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// Claim rewards
	validatorAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	_, err = k.Keeper.ClaimRewards(ctx, sender, validatorAddr)
	if err != nil {
		return nil, err
	}

	// Convert and delegate
	err = k.ConvertAndDelegate(ctx, sender, msg.ValidatorAddress, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Mint non transferrable share token
	shareDenom := types.NonVotingShareTokenDenom(msg.ValidatorAddress)
	k.bankKeeper.SetSendEnabled(ctx, shareDenom, false)

	coins := sdk.NewCoins(sdk.NewCoin(shareDenom, msg.Amount))

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	// Send non transferrable share token to sender
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coins)
	if err != nil {
		return nil, err
	}

	return &types.MsgNonVotingDelegateResponse{}, nil
}
