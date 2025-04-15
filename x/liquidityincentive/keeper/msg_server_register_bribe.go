package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
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

	// Check if epoch exists
	_, err = k.Epochs.Get(ctx, msg.EpochId)
	if err != nil {
		return nil, err
	}

	// Check if bribe already exists
	bribeKey := collections.Join(msg.EpochId, msg.PoolId)
	_, err = k.Bribes.Get(ctx, bribeKey)
	if err == nil {
		return nil, errorsmod.Wrap(types.ErrBribeAlreadyExists, "bribe already exists")
	}

	// Check if amount is valid
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return nil, errorsmod.Wrap(types.ErrInvalidBribe, "invalid bribe amount")
	}

	// Send coins from sender to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		sdkCtx,
		senderAddr,
		types.ModuleName,
		sdk.NewCoins(msg.Amount),
	); err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module")
	}

	// Create new bribe
	bribe := types.Bribe{
		Amount:        msg.Amount,
		ClaimedAmount: sdk.NewCoin(msg.Amount.Denom, math.ZeroInt()),
	}

	// Save bribe
	if err := k.Bribes.Set(ctx, bribeKey, bribe); err != nil {
		return nil, errorsmod.Wrap(err, "failed to save bribe")
	}

	// Emit event
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventRegisterBribe{
		Address: msg.Sender,
		EpochId: msg.EpochId,
		PoolId:  msg.PoolId,
		Amount:  msg.Amount,
	}); err != nil {
		return nil, err
	}

	return &types.MsgRegisterBribeResponse{}, nil
}
