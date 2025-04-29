package keeper

import (
	"context"
	"slices"

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
	if !found {
		return nil, types.ErrBribeNotFound
	}
	if err != nil {
		return nil, err
	}
	// Get bribe allocation
	allocation, err := k.GetBribeAllocation(ctx, senderAddr, bribe.EpochId, bribe.PoolId)
	if err != nil {
		return nil, err
	}

	// Check if bribe allocation is already claimed
	if slices.Contains(allocation.ClaimedBribeIds, bribe.Id) {
		return nil, types.ErrBribeAlreadyClaimed
	}

	// Get weight
	weight, err := math.LegacyNewDecFromStr(allocation.Weight)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid weight format")
	}

	// Calculate claim amount
	claimAmountDec := sdk.NewDecCoinsFromCoins(bribe.Amount...).MulDecTruncate(weight)
	claimAmount, _ := claimAmountDec.TruncateDecimal()

	if claimAmount.IsZero() {
		return nil, types.ErrNoBribeToClaim
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

	// Update bribe allocation (prevent double claiming)
	allocation.ClaimedBribeIds = append(allocation.ClaimedBribeIds, bribe.Id)
	if err := k.SetBribeAllocation(ctx, allocation); err != nil {
		return nil, errorsmod.Wrap(err, "failed to update bribe allocation")
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
