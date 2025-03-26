package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) NonVotingDelegate(ctx context.Context, msg *types.MsgNonVotingDelegate) (*types.MsgNonVotingDelegateResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// Validate amount
	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}
	if msg.Amount.Denom != feeDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "delegate amount denom must be equal to fee denom")
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
	err = k.ConvertAndDelegate(ctx, sender, msg.ValidatorAddress, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}

	// Mint non transferrable share token
	shareDenom := types.NonVotingShareTokenDenom(msg.ValidatorAddress)
	k.bankKeeper.SetSendEnabled(ctx, shareDenom, false)

	shareAmount, err := k.CalculateShareByAmount(ctx, msg.ValidatorAddress, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}

	coins := sdk.NewCoins(sdk.NewCoin(shareDenom, shareAmount))

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
