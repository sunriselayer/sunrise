package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k msgServer) RegisterBribe(ctx context.Context, msg *types.MsgRegisterBribe) (*types.MsgRegisterBribeResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	// Check if epoch exists
	epoch, found, err := k.Epochs.Get(ctx, msg.EpochId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrEpochNotFound, "epoch %d not found", msg.EpochId)
	}

	// Check if pool exists
	_, found, err = k.liquidityPoolKeeper.GetPool(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrPoolNotFound, "pool %d not found", msg.PoolId)
	}

	// Check if bribe amount is valid
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return nil, errorsmod.Wrap(types.ErrInvalidBribeAmount, "bribe amount must be valid and non-zero")
	}

	// Withdraw bribe amount from sender
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(sdkCtx, sender, types.ModuleName, msg.Amount); err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module")
	}

	// Save or update bribe
	key := collections.Join(msg.EpochId, msg.PoolId)
	existingBribe, found, err := k.Bribes.Get(ctx, key)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get existing bribe")
	}

	if found {
		// Add to existing bribe
		existingBribe.Amount = existingBribe.Amount.Add(msg.Amount)
		if err := k.Bribes.Set(ctx, key, existingBribe); err != nil {
			return nil, errorsmod.Wrap(err, "failed to update bribe")
		}
	} else {
		// Create new bribe
		bribe := types.Bribe{
			EpochId:       msg.EpochId,
			PoolId:        msg.PoolId,
			Amount:        msg.Amount,
			ClaimedAmount: sdk.NewCoin(msg.Amount.Denom, sdk.ZeroInt()),
		}

		if err := k.Bribes.Set(ctx, key, bribe); err != nil {
			return nil, errorsmod.Wrap(err, "failed to set bribe")
		}
	}

	// Emit event
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventRegisterBribe{
		Sender:  msg.Sender,
		EpochId: msg.EpochId,
		PoolId:  msg.PoolId,
		Amount:  msg.Amount,
	}); err != nil {
		return nil, err
	}

	return &types.MsgRegisterBribeResponse{}, nil
}
