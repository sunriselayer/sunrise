package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k msgServer) ClaimBribe(ctx context.Context, msg *types.MsgClaimBribe) (*types.MsgClaimBribeResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	senderAddr := sdk.AccAddress(sender)
	totalClaimed := sdk.NewCoins()

	// Check if epoch exists
	_, found, err := k.Epochs.Get(ctx, msg.EpochId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrEpochNotFound, "epoch %d not found", msg.EpochId)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Process bribes for each pool
	for _, poolId := range msg.PoolIds {
		// Check if bribe exists
		bribeKey := collections.Join(msg.EpochId, poolId)
		bribe, found, err := k.Bribes.Get(ctx, bribeKey)
		if err != nil {
			return nil, err
		}
		if !found {
			continue // No bribe for this pool
		}

		// Get unclaimed bribe
		unclaimedKey := collections.Join3(msg.Sender, msg.EpochId, poolId)
		unclaimed, found, err := k.UnclaimedBribes.Get(ctx, unclaimedKey)
		if err != nil {
			return nil, err
		}
		if !found {
			continue // No claim right
		}

		// Get weight
		weight, err := math.LegacyNewDecFromStr(unclaimed.Weight)
		if err != nil {
			return nil, errorsmod.Wrap(err, "invalid weight format")
		}

		// Calculate claim amount
		claimAmount := sdk.NewCoin(
			bribe.Amount.Denom,
			math.LegacyNewDecFromInt(bribe.Amount.Amount).Mul(weight).TruncateInt(),
		)

		if claimAmount.IsZero() {
			continue
		}

		// Send bribe
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			sdkCtx,
			types.ModuleName,
			senderAddr,
			sdk.NewCoins(claimAmount),
		); err != nil {
			return nil, errorsmod.Wrap(err, "failed to send coins from module")
		}

		// Update claimed amount
		bribe.ClaimedAmount = bribe.ClaimedAmount.Add(claimAmount)
		if err := k.Bribes.Set(ctx, bribeKey, bribe); err != nil {
			return nil, errorsmod.Wrap(err, "failed to update bribe claimed amount")
		}

		// Remove UnclaimedBribe (prevent double claiming)
		if err := k.UnclaimedBribes.Remove(ctx, unclaimedKey); err != nil {
			return nil, errorsmod.Wrap(err, "failed to remove unclaimed bribe")
		}

		totalClaimed = totalClaimed.Add(claimAmount)
	}

	if totalClaimed.IsZero() {
		return nil, errorsmod.Wrap(types.ErrNoBribesToClaim, "no bribes to claim")
	}

	// Emit event
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventClaimBribe{
		Sender:  msg.Sender,
		EpochId: msg.EpochId,
		PoolIds: msg.PoolIds,
		Amount:  totalClaimed,
	}); err != nil {
		return nil, err
	}

	return &types.MsgClaimBribeResponse{
		ClaimedAmount: totalClaimed,
	}, nil
}
