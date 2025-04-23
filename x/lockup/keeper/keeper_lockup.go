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

	// return error if the delegation amount is zero or if the base coins does not
	// exceed the desired delegation amount.
	if amount.IsZero() || balance.LT(amount) {
		return sdkerrors.ErrInvalidCoins.Wrap("delegation attempt with zero coins for staking denom or insufficient funds")
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

func (k Keeper) TrackUndelegation(ctx context.Context,
	owner sdk.AccAddress,
	id uint64,
	amount math.Int,
) error {
	lockup, err := k.GetLockupAccount(ctx, owner, id)
	if err != nil {
		return err
	}

	// return error if the delegation amount is zero or if the base coins does not
	// exceed the desired delegation amount.
	if amount.IsZero() {
		return sdkerrors.ErrInvalidCoins.Wrap("undelegation attempt with zero coins for staking denom")
	}

	delFree := lockup.DelegatedFree
	delLocking := lockup.DelegatedLocking

	// compute x and y per the specification, where:
	// X := min(DF, D)
	// Y := min(DV, D - X)
	x := math.MinInt(delFree, amount)
	y := math.MinInt(delLocking, amount.Sub(x))

	if !x.IsZero() {
		newDelFree := delFree.Sub(x)
		lockup.DelegatedFree = newDelFree
	}

	if !y.IsZero() {
		newDelLocking := delLocking.Sub(y)
		lockup.DelegatedLocking = newDelLocking
	}

	return k.SetLockupAccount(ctx, lockup)
}

// checkUnbondingEntriesMature iterates through all the unbonding entries and check if any of the entries are matured and handled.
func (k Keeper) CheckUnbondingEntriesMature(ctx context.Context, owner sdk.AccAddress, id uint64) error {
	lockup, err := k.GetLockupAccount(ctx, owner, id)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime().Unix()
	matureEntriesFound := false

	var keptEntries []*types.UnbondingEntry

	for _, entry := range lockup.UnbondEntries.Entries {
		if entry.EndTime > currentTime {
			keptEntries = append(keptEntries, entry)
		} else {
			matureEntriesFound = true
			err = k.TrackUndelegation(ctx, owner, id, entry.Amount)
			if err != nil {
				return err
			}
		}
	}

	if matureEntriesFound {
		lockup.UnbondEntries.Entries = keptEntries

		return k.SetLockupAccount(ctx, lockup)

	}

	return nil
}
