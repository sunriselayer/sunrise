package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k Keeper) GetLockupAccount(ctx context.Context, owner sdk.AccAddress, id uint64) (types.LockupAccount, error) {
	lockupAccount, err := k.LockupAccounts.Get(ctx, collections.Join([]byte(owner), id))
	if err != nil {
		return types.LockupAccount{}, err
	}

	return lockupAccount, nil
}

func (k Keeper) SetLockupAccount(ctx context.Context, lockupAccount types.LockupAccount) error {
	address, err := k.addressCodec.StringToBytes(lockupAccount.Owner)
	if err != nil {
		return err
	}

	return k.LockupAccounts.Set(ctx, collections.Join(address, lockupAccount.Id), lockupAccount)
}

func (k Keeper) GetAllLockupAccounts(ctx context.Context) (list []types.LockupAccount, err error) {
	err = k.LockupAccounts.Walk(ctx, nil, func(key collections.Pair[[]byte, uint64], value types.LockupAccount) (stop bool, err error) {
		list = append(list, value)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (k Keeper) GetLockupAccountsByOwner(ctx context.Context, owner sdk.AccAddress) (list []types.LockupAccount, err error) {
	err = k.LockupAccounts.Walk(
		ctx,
		collections.NewPrefixedPairRange[[]byte, uint64](owner),
		func(key collections.Pair[[]byte, uint64], value types.LockupAccount) (stop bool, err error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
