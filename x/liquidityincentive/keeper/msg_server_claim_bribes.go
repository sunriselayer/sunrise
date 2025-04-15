package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k msgServer) ClaimBribes(ctx context.Context, msg *types.MsgClaimBribes) (*types.MsgClaimBribesResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	senderAddr := sdk.AccAddress(sender)
	totalClaimed := sdk.NewCoins()

	// Check if epoch exists
	_, err = k.Epochs.Get(ctx, msg.EpochId)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check if bribe exists
	bribeKey := collections.Join(msg.EpochId, msg.PoolId)
	bribe, err := k.Bribes.Get(ctx, bribeKey)
	if err != nil {
		return nil, err
	}

	// Get unclaimed bribe
	unclaimedKey := collections.Join3(senderAddr, msg.EpochId, msg.PoolId)
	unclaimed, err := k.UnclaimedBribes.Get(ctx, unclaimedKey)
	if err != nil {
		return nil, err
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
		return nil, errorsmod.Wrap(types.ErrNoBribesToClaim, "no bribes to claim")
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

	// Emit event
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventClaimBribes{
		Address:       msg.Sender,
		ClaimedBribes: totalClaimed,
	}); err != nil {
		return nil, err
	}

	return &types.MsgClaimBribesResponse{
		ClaimedBribes: totalClaimed,
	}, nil
}
