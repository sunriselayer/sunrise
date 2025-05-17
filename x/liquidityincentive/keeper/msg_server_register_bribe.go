package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k msgServer) RegisterBribe(ctx context.Context, msg *types.MsgRegisterBribe) (*types.MsgRegisterBribeResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	senderAddr := sdk.AccAddress(sender)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check if epoch is in the future
	currentEpoch, found, err := k.GetLastEpoch(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get current epoch")
	}

	if found && msg.EpochId < currentEpoch.Id {
		return nil, errorsmod.Wrap(types.ErrInvalidBribe, "epoch is in the past")
	}

	// Check if amount is valid
	err = msg.Amount.Validate()
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid bribe amount")
	}

	if msg.Amount.IsZero() {
		return nil, errorsmod.Wrap(types.ErrInvalidBribe, "amount cannot be zero")
	}

	// Check if send enabled coins
	if err := k.bankKeeper.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, errorsmod.Wrap(err, "failed to check if send enabled coins")
	}

	// Check if pool exists before any bank keeper calls
	_, found, err = k.liquidityPoolKeeper.GetPool(ctx, msg.PoolId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get pool")
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrInvalidBribe, "pool %d not found", msg.PoolId)
	}

	// Send coins from sender to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		sdkCtx,
		senderAddr,
		types.BribeAccount,
		msg.Amount,
	); err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module")
	}

	// Create new bribe
	bribe := types.Bribe{
		EpochId:       msg.EpochId,
		PoolId:        msg.PoolId,
		Address:       msg.Sender,
		Amount:        msg.Amount,
		ClaimedAmount: sdk.NewCoins(),
	}

	// Save bribe
	id, err := k.AppendBribe(ctx, bribe)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to save bribe")
	}

	// Emit event
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventRegisterBribe{
		Id:      id,
		EpochId: msg.EpochId,
		PoolId:  msg.PoolId,
		Address: msg.Sender,
		Amount:  msg.Amount,
	}); err != nil {
		return nil, err
	}

	return &types.MsgRegisterBribeResponse{}, nil
}
