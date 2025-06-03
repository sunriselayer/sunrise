package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k Keeper) LockupAccountAddress(owner sdk.AccAddress, id uint64) sdk.AccAddress {
	seed := k.makeAddressSeed(owner, id)
	return authtypes.NewModuleAddress(seed)
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
			// Owner not found, start from ID 1
			currentID = 1
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

func (k Keeper) InitLockupAccountFromMsg(ctx context.Context, msg *types.MsgInitLockupAccount) error {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return errorsmod.Wrap(err, "invalid sender address")
	}

	owner, err := k.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return errorsmod.Wrap(err, "invalid owner address")
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return err
	}

	if msg.Amount.Denom != feeDenom {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount denom must be fee denom")
	}

	// Get the current ID for the owner and increment the counter
	id, _, err := k.GetAndIncrementNextLockupAccountID(ctx, owner)
	if err != nil {
		return errorsmod.Wrap(err, "failed to get next lockup account ID")
	}

	// Generate the seed and the lockup account address
	lockupAcc := k.LockupAccountAddress(owner, id)

	err = k.bankKeeper.SendCoins(ctx, sender, lockupAcc, sdk.NewCoins(msg.Amount))
	if err != nil {
		return errorsmod.Wrap(err, "failed to send coins")
	}

	lockupAccount := types.LockupAccount{
		Address:           lockupAcc.String(),
		Owner:             msg.Owner,
		Id:                id,
		StartTime:         msg.StartTime,
		EndTime:           msg.EndTime,
		OriginalLocking:   msg.Amount.Amount,
		AdditionalLocking: math.ZeroInt(),
		DelegatedFree:     math.ZeroInt(),
		DelegatedLocking:  math.ZeroInt(),
		UnbondEntries:     &types.UnbondingEntries{},
	}

	err = k.SetLockupAccount(ctx, lockupAccount)
	if err != nil {
		return err
	}

	return nil
}
