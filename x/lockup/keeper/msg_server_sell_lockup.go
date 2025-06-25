package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

// SellLockupAccount handles MsgSellLockupAccount.
// It allows a lockup account owner to list their account for sale.
func (k msgServer) SellLockupAccount(ctx context.Context, msg *types.MsgSellLockupAccount) (*types.MsgSellLockupAccountResponse, error) {
	seller, err := k.accountKeeper.AddressCodec().StringToBytes(msg.Owner)
	if err != nil {
		return nil, err
	}

	// Check if the lockup account exists
	if _, err := k.LockupAccounts.Get(ctx, collections.Join(seller, msg.LockupAccountId)); err != nil {
		return nil, err
	}

	// Check if the lockup account is already listed for sale
	if _, err := k.Listings.Get(ctx, collections.Join(seller, msg.LockupAccountId)); err == nil {
		return nil, types.ErrAlreadyListed
	}

	// List the lockup account for sale
	if err := k.Listings.Set(ctx, collections.Join(seller, msg.LockupAccountId), msg.Price); err != nil {
		return nil, err
	}

	return &types.MsgSellLockupAccountResponse{}, nil
}
