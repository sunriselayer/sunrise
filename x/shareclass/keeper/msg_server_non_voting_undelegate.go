package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) NonVotingUndelegate(ctx context.Context, msg *types.MsgNonVotingUndelegate) (*types.MsgNonVotingUndelegateResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
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

	validatorAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	err = msg.Amount.Validate()
	if err != nil {
		return nil, err
	}

	output, rewards, completionTime, err := k.Undelegate(ctx, sender, recipient, validatorAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgNonVotingUndelegateResponse{
		CompletionTime: completionTime,
		Amount:         output,
		Rewards:        rewards,
	}, nil
}
