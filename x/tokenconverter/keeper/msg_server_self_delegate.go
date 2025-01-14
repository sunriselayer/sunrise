package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "cosmossdk.io/x/staking/types"

	selfdelegationproxy "github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy"
	selfdelegationproxytypes "github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy/v1"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k msgServer) SelfDelegate(ctx context.Context, msg *types.MsgSelfDelegate) (*types.MsgSelfDelegateResponse, error) {
	delegator, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, err
	}

	rootOwnerAcc, rootOwnerVal, err := k.getRootOwner(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}

	// Check proxy account existence
	has, err := k.SelfDelegationProxies.Has(ctx, delegator)
	if err != nil {
		return nil, err
	}
	var proxyAddrBytes []byte
	if has {
		proxyAddrBytes, err = k.SelfDelegationProxies.Get(ctx, delegator)
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
				Owner:     msg.Sender,
				RootOwner: rootOwnerAcc,
			},
			sdk.NewCoins(), // Do not use to unify the case of already existing proxy account
			[]byte{},
		)
		if err != nil {
			return nil, err
		}
		err = k.SelfDelegationProxies.Set(ctx, delegator, proxyAddrBytes)
		if err != nil {
			return nil, err
		}
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoins(ctx, delegator, proxyAddrBytes, sdk.NewCoins(sdk.NewCoin(params.FeeDenom, msg.Amount)))
	if err != nil {
		return nil, err
	}

	// ConvertReverse
	err = k.ConvertReverse(ctx, msg.Amount, proxyAddrBytes)
	if err != nil {
		return nil, err
	}

	// Delegate from proxy account
	proxyAddr, err := k.addressCodec.BytesToString(proxyAddrBytes)
	if err != nil {
		return nil, err
	}
	_, err = k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: proxyAddr,
		ValidatorAddress: rootOwnerVal,
		Amount:           sdk.NewCoin(params.BondDenom, msg.Amount),
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgSelfDelegateResponse{}, nil
}
