package keeper

import (
	"context"

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

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check if bribe exists
	bribe, found, err := k.GetBribe(ctx, msg.BribeId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, types.ErrBribeNotFound
	}
	// Get unclaimed bribe
	unclaimed, err := k.GetUnclaimedBribe(ctx, senderAddr, msg.BribeId)
	if err != nil {
		return nil, err
	}

	// Get weight
	weight, err := math.LegacyNewDecFromStr(unclaimed.Weight)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid weight format")
	}

	// Calculate claim amount
	claimAmountDec := sdk.NewDecCoinsFromCoins(bribe.Amount...).MulDecTruncate(weight)
	claimAmount, _ := claimAmountDec.TruncateDecimal()

	if claimAmount.IsZero() {
		return nil, errorsmod.Wrap(types.ErrNoBribesToClaim, "no bribes to claim")
	}

	// Send bribe
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		sdkCtx,
		types.BribeAccount,
		senderAddr,
		claimAmount,
	); err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins from module")
	}

	// Update claimed amount
	bribe.ClaimedAmount = bribe.ClaimedAmount.Add(claimAmount...)
	if err := k.SetBribe(ctx, bribe); err != nil {
		return nil, errorsmod.Wrap(err, "failed to update bribe claimed amount")
	}

	// Remove UnclaimedBribe (prevent double claiming)
	if err := k.RemoveUnclaimedBribe(ctx, unclaimed); err != nil {
		return nil, errorsmod.Wrap(err, "failed to remove unclaimed bribe")
	}

	// Emit event
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventClaimBribes{
		Address: msg.Sender,
		Amount:  claimAmount,
	}); err != nil {
		return nil, err
	}

	return &types.MsgClaimBribesResponse{
		Amount: claimAmount,
	}, nil
}
