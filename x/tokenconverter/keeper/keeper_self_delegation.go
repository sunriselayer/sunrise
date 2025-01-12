package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/x/accounts/defaults/lockup"
	lockuptypes "cosmossdk.io/x/accounts/defaults/lockup/v1"
	stakingtypes "cosmossdk.io/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	accounttypes "cosmossdk.io/x/accounts/v1"

	selfdelegationproxy "github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy"
	selfdelegationproxytypes "github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy/v1"
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

func (k Keeper) DelegateOrSelfDelegate(
	ctx context.Context,
	msg *stakingtypes.MsgDelegate,
	originalFunc func(ctx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error),
) (*stakingtypes.MsgDelegateResponse, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Denom == params.FeeDenom {
		delegator, err := k.addressCodec.StringToBytes(msg.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		rootOwnerAcc, rootOwnerVal, err := k.getRootOwner(ctx, msg.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		if rootOwnerVal != msg.ValidatorAddress {
			return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "delegation with denom %s is only for self delegation", params.FeeDenom)
		}

		// Check proxy account existence
		has, err := k.SelfDelegationProxy.Has(ctx, delegator)
		if err != nil {
			return nil, err
		}
		var proxyAddrBytes []byte
		if has {
			proxyAddrBytes, err = k.SelfDelegationProxy.Get(ctx, delegator)
			if err != nil {
				return nil, err
			}
		} else {
			// Create proxy account
			_, proxyAddrBytes, err = k.accountsKeeper.Init(
				ctx,
				selfdelegationproxy.SELF_DELEGATION_PROXY_ACCOUNT,
				delegator, // Must be delegator, not owner
				&selfdelegationproxytypes.MsgInit{
					Owner:     msg.DelegatorAddress,
					RootOwner: rootOwnerAcc,
				},
				sdk.NewCoins(), // Do not use to unify the case of already existing proxy account
				[]byte{},
			)
			if err != nil {
				return nil, err
			}
			err = k.SelfDelegationProxy.Set(ctx, delegator, proxyAddrBytes)
			if err != nil {
				return nil, err
			}
		}

		err = k.bankKeeper.SendCoins(ctx, delegator, proxyAddrBytes, sdk.NewCoins(msg.Amount))
		if err != nil {
			return nil, err
		}

		// ConvertReverse
		err = k.ConvertReverse(ctx, msg.Amount.Amount, proxyAddrBytes)
		if err != nil {
			return nil, err
		}

		// Delegate from proxy account
		proxyAddr, err := k.addressCodec.BytesToString(proxyAddrBytes)
		if err != nil {
			return nil, err
		}
		res, err := k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgDelegate{
			DelegatorAddress: proxyAddr,
			ValidatorAddress: rootOwnerVal,
			Amount:           sdk.NewCoin(params.BondDenom, msg.Amount.Amount),
		})

		return res.(*stakingtypes.MsgDelegateResponse), err
	}

	return originalFunc(ctx, msg)
}

func (k Keeper) UndelegateOrSelfUndelegate(
	ctx context.Context,
	msg *stakingtypes.MsgUndelegate,
	originalFunc func(ctx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error),
) (*stakingtypes.MsgUndelegateResponse, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Denom == params.FeeDenom {

		delegator, err := k.addressCodec.StringToBytes(msg.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		_, rootOwnerVal, err := k.getRootOwner(ctx, msg.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		if rootOwnerVal != msg.ValidatorAddress {
			return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "delegation with denom %s is only for self delegation", params.FeeDenom)
		}

		proxyAddrBytes, err := k.SelfDelegationProxy.Get(ctx, delegator)
		if err != nil {
			return nil, err
		}
		proxyAddr, err := k.addressCodec.BytesToString(proxyAddrBytes)
		if err != nil {
			return nil, err
		}

		res, err := k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgUndelegate{
			DelegatorAddress: proxyAddr,
			ValidatorAddress: rootOwnerVal,
			Amount:           sdk.NewCoin(params.BondDenom, msg.Amount.Amount),
		})

		return res.(*stakingtypes.MsgUndelegateResponse), err
	}

	return originalFunc(ctx, msg)
}
