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

	// Get LST supply before burning
	shareDenom := types.NonVotingShareTokenDenom(msg.ValidatorAddress)

	// Send non transferrable share token to module
	coins := sdk.NewCoins(sdk.NewCoin(shareDenom, msg.Share))

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

	// Calculate unbonding amount
	totalShare := k.GetTotalShare(ctx, msg.ValidatorAddress)
	totalStaked, err := k.GetTotalStakedAmount(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	outputAmount, err := types.CalculateUndelegationOutputAmount(msg.Share, totalShare, totalStaked)
	if err != nil {
		return nil, err
	}

	// Undelegate
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}
	output := sdk.NewCoin(bondDenom, outputAmount)

	res, err := k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgUndelegate{
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
