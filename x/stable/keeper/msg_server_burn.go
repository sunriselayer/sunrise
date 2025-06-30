package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/stable/types"
)

func (k msgServer) Burn(ctx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "invalid sender address: %s", msg.Sender)
	}

	if err := sdk.ValidateDenom(msg.OutputDenom); err != nil {
		return nil, errorsmod.Wrap(err, "invalid output denom")
	}

	if !msg.Amount.IsPositive() {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "burn amount must be positive")
	}

	returnedCoins, err := k.Keeper.Burn(ctx, sender, msg.Amount, msg.OutputDenom)
	if err != nil {
		return nil, err
	}

	return &types.MsgBurnResponse{
		Amount: returnedCoins,
	}, nil
}
