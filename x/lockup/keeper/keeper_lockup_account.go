package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) LockupAccountAddress(owner sdk.AccAddress, id uint64) sdk.AccAddress {
	seed := k.makeAddressSeed(owner, id)
	return k.accountKeeper.GetModuleAddress(seed)
}

func (k Keeper) makeAddressSeed(owner sdk.AccAddress, id uint64) string {
	return fmt.Sprintf("lockup/%s/%d", owner.String(), id)
}

// GetAndIncrementNextLockupAccountID retrieves the current ID for the owner,
// increments the counter, and returns both the current and the next ID.
// If the owner does not exist, it returns 0 as the current ID and 1 as the next ID.
func (k Keeper) GetAndIncrementNextLockupAccountID(ctx context.Context, owner sdk.AccAddress) (currentID uint64, nextID uint64, err error) {
	currentID, err = k.NextLockupAccountId.Get(ctx, owner)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Owner not found, start from ID 0
			currentID = 0
		} else {
			// Other error occurred
			return 0, 0, fmt.Errorf("failed to get next lockup account id for owner %s: %w", owner.String(), err)
		}
	}

	nextID = currentID + 1

	// Set the next ID for the owner
	err = k.NextLockupAccountId.Set(ctx, owner, nextID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to set next lockup account id %d for owner %s: %w", nextID, owner.String(), err)
	}

	return currentID, nextID, nil
}
