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

func (k Keeper) SelfDelegate(ctx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// If the amount is in the fee denom
	if msg.Amount.Denom == params.FeeDenom {
		res, err := k.Environment.QueryRouterService.Invoke(ctx, &accounttypes.AccountTypeRequest{
			Address: msg.DelegatorAddress,
		})
		if err != nil {
			return nil, err
		}

		delegator, err := k.addressCodec.StringToBytes(msg.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		var rootOwner []byte
		switch res.(*accounttypes.AccountTypeResponse).AccountType {
		// Case of lockup accounts
		case lockup.CONTINUOUS_LOCKING_ACCOUNT,
			lockup.DELAYED_LOCKING_ACCOUNT,
			lockup.PERIODIC_LOCKING_ACCOUNT,
			lockup.PERMANENT_LOCKING_ACCOUNT:

			res, err := k.accountKeeper.Query(ctx, delegator, &lockuptypes.QueryLockupAccountInfoRequest{})
			if err != nil {
				return nil, err
			}
			rootOwner, err = k.addressCodec.StringToBytes(res.(*lockuptypes.QueryLockupAccountInfoResponse).Owner)
			if err != nil {
				return nil, err
			}
		default:
			rootOwner = delegator
		}
		// Convert to val address
		rootOwnerValAddress, err := k.validatorAddressCodec.BytesToString(rootOwner)
		if err != nil {
			return nil, err
		}

		if rootOwnerValAddress != msg.ValidatorAddress {
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
			rootOwnerString, err := k.addressCodec.BytesToString(rootOwner)
			if err != nil {
				return nil, err
			}

			// Create proxy account
			_, proxyAddrBytes, err = k.accountKeeper.Init(
				ctx,
				selfdelegationproxy.SELF_DELEGATION_PROXY_ACCOUNT,
				delegator, // Must be delegator, not owner
				&selfdelegationproxytypes.MsgInit{
					Owner:     msg.DelegatorAddress,
					RootOwner: rootOwnerString,
				},
				sdk.NewCoins(msg.Amount),
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
		res, err = k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgDelegate{
			DelegatorAddress: proxyAddr,
			ValidatorAddress: rootOwnerValAddress,
			Amount:           sdk.NewCoin(params.BondDenom, msg.Amount.Amount),
		})

		return res.(*stakingtypes.MsgDelegateResponse), err
	}

	return nil, nil
}
