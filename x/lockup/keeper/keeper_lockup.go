package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k Keeper) TrackDelegation(
	ctx context.Context,
	owner sdk.AccAddress,
	id uint64,
	balance math.Int,
	locked math.Int,
	amount math.Int,
) error {
	lockup, err := k.GetLockupAccount(ctx, owner, id)
	if err != nil {
		return err
	}

	// return error if the delegation amount is zero
	if amount.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrap("delegation amount cannot be zero")
	}

	// return error if the delegation amount exceeds total balance
	if balance.LT(amount) {
		return sdkerrors.ErrInsufficientFunds.Wrapf("total balance is less than delegation amount: %s < %s", balance, amount)
	}

	delLocking := lockup.DelegatedLocking
	delFree := lockup.DelegatedFree

	// compute x and y per the specification, where:
	// X := min(max(V - DV, 0), D)
	// Y := D - X
	x := math.MinInt(math.MaxInt(locked.Sub(delLocking), math.ZeroInt()), amount)
	y := amount.Sub(x)

	if !x.IsZero() {
		newDelLocking := delLocking.Add(x)
		lockup.DelegatedLocking = newDelLocking
	}

	if !y.IsZero() {
		newDelFree := delFree.Add(y)
		lockup.DelegatedFree = newDelFree
	}

	return k.SetLockupAccount(ctx, lockup)
}

// CheckUnbondingEntriesMature iterates through all the unbonding entries, handles matured entries,
// and updates the lockup account state in a single operation.
func (k Keeper) CheckUnbondingEntriesMature(ctx context.Context, owner sdk.AccAddress, id uint64) error {
	lockup, err := k.GetLockupAccount(ctx, owner, id)
	if err != nil {
		return err
	}

	if lockup.UnbondEntries == nil || len(lockup.UnbondEntries.Entries) == 0 {
		return nil
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime().Unix()

	var keptEntries []*types.UnbondingEntry
	totalMatureAmount := math.ZeroInt()

	for _, entry := range lockup.UnbondEntries.Entries {
		if entry.EndTime > currentTime {
			keptEntries = append(keptEntries, entry)
		} else {
			totalMatureAmount = totalMatureAmount.Add(entry.Amount)
		}
	}

	// If no entries have matured, there's nothing to do.
	if totalMatureAmount.IsZero() {
		// still need to update the entries if some have matured
		if len(keptEntries) < len(lockup.UnbondEntries.Entries) {
			lockup.UnbondEntries.Entries = keptEntries
			return k.SetLockupAccount(ctx, lockup)
		}
		return nil
	}

	// Apply the logic from TrackUndelegation directly to the in-memory lockup object.
	// This avoids multiple store reads and writes.
	delFree := lockup.DelegatedFree
	delLocking := lockup.DelegatedLocking

	// X := min(DF, D)
	// Y := min(DV, D - X)
	x := math.MinInt(delFree, totalMatureAmount)
	y := math.MinInt(delLocking, totalMatureAmount.Sub(x))

	if !x.IsZero() {
		lockup.DelegatedFree = delFree.Sub(x)
	}
	if !y.IsZero() {
		lockup.DelegatedLocking = delLocking.Sub(y)
	}

	// Update the unbonding entries list.
	lockup.UnbondEntries.Entries = keptEntries

	// Set the updated lockup account once.
	return k.SetLockupAccount(ctx, lockup)
}

func (k Keeper) AddAdditionalLockup(ctx context.Context, owner sdk.AccAddress, id uint64, amount math.Int) error {
	lockup, err := k.GetLockupAccount(ctx, owner, id)
	if err != nil {
		return err
	}

	lockup.AdditionalLocking = lockup.AdditionalLocking.Add(amount)

	return k.SetLockupAccount(ctx, lockup)
}
