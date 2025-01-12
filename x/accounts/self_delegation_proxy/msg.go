package selfdelegationproxy

import (
	"bytes"
	"context"

	"cosmossdk.io/x/accounts/accountstd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	banktypes "cosmossdk.io/x/bank/types"
	distrtypes "cosmossdk.io/x/distribution/types"
	tokenconvertertypes "github.com/sunriselayer/sunrise/x/tokenconverter/types"

	v1 "github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy/v1"
)

func (a SelfDelegationProxyAccount) Init(ctx context.Context, msg *v1.MsgInit) (*v1.MsgInitResponse, error) {
	owner, err := a.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, err
	}
	err = a.Owner.Set(ctx, owner)
	if err != nil {
		return nil, err
	}

	rootOwner, err := a.addressCodec.StringToBytes(msg.RootOwner)
	if err != nil {
		return nil, err
	}
	err = a.RootOwner.Set(ctx, rootOwner)
	if err != nil {
		return nil, err
	}

	return &v1.MsgInitResponse{}, nil
}

func (a SelfDelegationProxyAccount) WithdrawReward(ctx context.Context, msg *v1.MsgWithdrawReward) (*v1.MsgWithdrawRewardResponse, error) {
	err := a.checkSender(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}
	whoami := accountstd.Whoami(ctx)
	delegatorAddress, err := a.addressCodec.BytesToString(whoami)
	if err != nil {
		return nil, err
	}

	msgWithdraw := &distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
	}
	res, err := accountstd.ExecModule[*distrtypes.MsgWithdrawDelegatorRewardResponse](ctx, msgWithdraw)
	if err != nil {
		return nil, err
	}

	return &v1.MsgWithdrawRewardResponse{
		Amount: res.Amount,
	}, nil
}

func (a SelfDelegationProxyAccount) Send(ctx context.Context, msg *v1.MsgSend) (*v1.MsgSendResponse, error) {
	err := a.checkSender(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}
	whoami := accountstd.Whoami(ctx)
	fromAddress, err := a.addressCodec.BytesToString(whoami)
	if err != nil {
		return nil, err
	}

	msgSend := &banktypes.MsgSend{
		FromAddress: fromAddress,
		ToAddress:   msg.ToAddress,
		Amount:      msg.Amount,
	}
	_, err = accountstd.ExecModule[*banktypes.MsgSendResponse](ctx, msgSend)
	if err != nil {
		return nil, err
	}

	return &v1.MsgSendResponse{}, nil
}

func (a SelfDelegationProxyAccount) WithdrawUnbonded(ctx context.Context, msg *v1.MsgWithdrawUnbonded) (*v1.MsgWithdrawUnbondedResponse, error) {
	err := a.checkSender(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}
	whoami := accountstd.Whoami(ctx)
	fromAddress, err := a.addressCodec.BytesToString(whoami)
	if err != nil {
		return nil, err
	}

	msgConvert := &tokenconvertertypes.MsgConvert{
		Sender: fromAddress,
		Amount: msg.Amount,
	}
	_, err = accountstd.ExecModule[*tokenconvertertypes.MsgConvertResponse](ctx, msgConvert)
	if err != nil {
		return nil, err
	}

	owner, err := a.Owner.Get(ctx)
	if err != nil {
		return nil, err
	}
	toAddress, err := a.addressCodec.BytesToString(owner)
	if err != nil {
		return nil, err
	}

	feeDenom, err := getFeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	msgSend := &banktypes.MsgSend{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      sdk.NewCoins(sdk.NewCoin(feeDenom, msg.Amount)),
	}
	_, err = accountstd.ExecModule[*banktypes.MsgSendResponse](ctx, msgSend)
	if err != nil {
		return nil, err
	}

	return &v1.MsgWithdrawUnbondedResponse{}, nil
}

func getBondDenom(ctx context.Context) (string, error) {
	params, err := accountstd.QueryModule[*tokenconvertertypes.QueryParamsResponse](ctx, &tokenconvertertypes.QueryParamsRequest{})
	if err != nil {
		return "", err
	}

	return params.Params.BondDenom, nil
}

func getFeeDenom(ctx context.Context) (string, error) {
	params, err := accountstd.QueryModule[*tokenconvertertypes.QueryParamsResponse](ctx, &tokenconvertertypes.QueryParamsRequest{})
	if err != nil {
		return "", err
	}

	return params.Params.FeeDenom, nil
}

func (a SelfDelegationProxyAccount) checkSender(ctx context.Context, sender string) error {
	rootOwner, err := a.RootOwner.Get(ctx)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err.Error())
	}
	senderBytes, err := a.addressCodec.StringToBytes(sender)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err.Error())
	}
	if !bytes.Equal(rootOwner, senderBytes) {
		return sdkerrors.ErrUnauthorized
	}

	return nil
}
