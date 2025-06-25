package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

// CancelSellLockupAccount handles MsgCancelSellLockupAccount.
// It allows a lockup account owner to cancel their sale listing.
func (k msgServer) CancelSellLockupAccount(ctx context.Context, msg *types.MsgCancelSellLockupAccount) (*types.MsgCancelSellLockupAccountResponse, error) {
	seller, err := k.accountKeeper.AddressCodec().StringToBytes(msg.Owner)
	if err != nil {
		return nil, err
	}

	// Check if the lockup account is listed for sale
	if _, err := k.Listings.Get(ctx, collections.Join(seller, msg.LockupAccountId)); err != nil {
		return nil, types.ErrNotListed
	}

	// Remove the listing
	if err := k.Listings.Remove(ctx, collections.Join(seller, msg.LockupAccountId)); err != nil {
		return nil, err
	}

	return &types.MsgCancelSellLockupAccountResponse{}, nil
}
