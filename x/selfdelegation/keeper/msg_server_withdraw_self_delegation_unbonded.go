package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/selfdelegation/types"
)

func (k msgServer) WithdrawSelfDelegationUnbonded(ctx context.Context, msg *types.MsgWithdrawSelfDelegationUnbonded) (*types.MsgWithdrawSelfDelegationUnbondedResponse, error) {
	delegator, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, err
	}

	proxyAddrBytes, err := k.SelfDelegationProxies.Get(ctx, delegator)
	if err != nil {
		return nil, err
	}

	err = k.tokenConverterKeeper.Convert(ctx, msg.Amount, proxyAddrBytes)
	if err != nil {
		return nil, err
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoins(ctx, proxyAddrBytes, delegator, sdk.NewCoins(sdk.NewCoin(feeDenom, msg.Amount)))
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawSelfDelegationUnbondedResponse{}, nil
}
