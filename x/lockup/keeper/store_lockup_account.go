package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k Keeper) GetLockupAccount(ctx context.Context, address sdk.AccAddress) (types.LockupAccount, error) {
	lockupAccount, err := k.LockupAccounts.Get(ctx, address)
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

	return k.LockupAccounts.Set(ctx, address, lockupAccount)
}
