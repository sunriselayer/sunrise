package keeper

import (
	"context"

	lockuptypes "cosmossdk.io/x/accounts/defaults/lockup/v1"
	accounttypes "cosmossdk.io/x/accounts/v1"

	lockup "github.com/sunriselayer/sunrise/x/accounts/self_delegatable_lockup"
)

func (k Keeper) getRootOwner(ctx context.Context, delegator string) (accAddress string, valAddress string, err error) {
	res, err := k.Environment.QueryRouterService.Invoke(ctx, &accounttypes.AccountTypeRequest{
		Address: delegator,
	})
	if err != nil {
		return
	}

	delegatorBytes, err := k.addressCodec.StringToBytes(delegator)
	if err != nil {
		return
	}

	var rootOwner []byte
	switch res.(*accounttypes.AccountTypeResponse).AccountType {
	// Case of lockup accounts
	case lockup.CONTINUOUS_LOCKING_ACCOUNT,
		lockup.DELAYED_LOCKING_ACCOUNT,
		lockup.PERIODIC_LOCKING_ACCOUNT,
		lockup.PERMANENT_LOCKING_ACCOUNT:

		res, err = k.accountsKeeper.Query(ctx, delegatorBytes, &lockuptypes.QueryLockupAccountInfoRequest{})
		if err != nil {
			return
		}
		rootOwner, err = k.addressCodec.StringToBytes(res.(*lockuptypes.QueryLockupAccountInfoResponse).Owner)
		if err != nil {
			return
		}
	default:
		rootOwner = delegatorBytes
	}
	// Convert to acc address
	accAddress, err = k.addressCodec.BytesToString(rootOwner)
	if err != nil {
		return
	}

	// Convert to val address
	valAddress, err = k.validatorAddressCodec.BytesToString(rootOwner)
	if err != nil {
		return
	}

	return
}
