package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stakingtypes "cosmossdk.io/x/staking/types"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) NonVotingUndelegate(ctx context.Context, msg *types.MsgNonVotingUndelegate) (*types.MsgNonVotingUndelegateResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	// Validate amount
	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}
	if msg.Amount.Denom != feeDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "undelegate amount denom must be equal to fee denom")
	}
	if !msg.Amount.IsPositive() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "undelegate amount must be positive")
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

	// Calculate unbonding share
	shareDenom := types.NonVotingShareTokenDenom(msg.ValidatorAddress)
	unbondingShare, err := k.CalculateShareByAmount(ctx, msg.ValidatorAddress, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}

	// Send non transferrable share token to module
	coins := sdk.NewCoins(sdk.NewCoin(shareDenom, unbondingShare))

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	// Burn non transferrable share token
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	err = k.bankKeeper.BurnCoins(ctx, moduleAddr, coins)
	if err != nil {
		return nil, err
	}

	// Undelegate
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}
	output := sdk.NewCoin(bondDenom, msg.Amount.Amount)

	res, err := k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgUndelegate{
		DelegatorAddress: moduleAddr.String(),
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
	if undelegateResponse == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "undelegate response is nil")
	}

	// Set recipient
	var recipient sdk.AccAddress
	if msg.Recipient == "" {
		recipient = sender
	} else {
		recipient, err = k.addressCodec.StringToBytes(msg.Recipient)
		if err != nil {
			return nil, errorsmod.Wrap(err, "invalid recipient address")
		}
	}

	// Append Unstaking state
	_, err = k.AppendUnbonding(ctx, types.Unbonding{
		Address:        recipient.String(),
		CompletionTime: undelegateResponse.CompletionTime,
		Amount:         output,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgNonVotingUndelegateResponse{
		CompletionTime: undelegateResponse.CompletionTime,
		Amount:         output,
	}, nil
}
