package selfdelegationproxy

import (
	"context"

	"github.com/cosmos/gogoproto/proto"

	"cosmossdk.io/math"
	"cosmossdk.io/x/accounts/accountstd"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	banktypes "cosmossdk.io/x/bank/types"
	distrtypes "cosmossdk.io/x/distribution/types"
	tokenconvertertypes "github.com/sunriselayer/sunrise/x/tokenconverter/types"

	v1 "github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy/v1"
)

func (a SelfDelegationProxyAccount) Init(ctx context.Context, msg *v1.MsgInit) (*v1.MsgInitResponse, error) {
	// Save parent account

	return &v1.MsgInitResponse{}, nil
}

func (a SelfDelegationProxyAccount) Undelegate(ctx context.Context, msg *v1.MsgUndelegate) (*v1.MsgUndelegateResponse, error) {

	return &v1.MsgUndelegateResponse{}, nil
}

func (a SelfDelegationProxyAccount) CancelUnbonding(ctx context.Context, msg *v1.MsgCancelUnbonding) (*v1.MsgCancelUnbondingResponse, error) {

	return &v1.MsgCancelUnbondingResponse{}, nil
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
	responses, err := executeMsg(ctx, msgWithdraw)
	if err != nil {
		return nil, err
	}

	return &v1.MsgWithdrawRewardResponse{}, nil
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

	hs := a.headerService.HeaderInfo(ctx)

	if err := msg.Amount.Validate(); err != nil {
		return nil, err
	}

	msgSend := &banktypes.MsgSend{
		FromAddress: fromAddress,
		ToAddress:   msg.ToAddress,
		Amount:      msg.Amount,
	}
	resp, err := executeMsg(ctx, msgSend)
	if err != nil {
		return nil, err
	}

	return &v1.MsgSendResponse{}, nil
}

func (a SelfDelegationProxyAccount) WithdrawUnbonded(ctx context.Context, msg *v1.MsgWithdrawUnbonded) (*v1.MsgWithdrawUnbondedResponse, error) {
	msgConvert := &tokenconvertertypes.MsgConvert{
		Amount: msg.Amount,
	}
	resp, err := executeMsg(ctx, msgConvert)
	if err != nil {
		return nil, err
	}

	msgSend := &banktypes.MsgSend{}
	resp, err = executeMsg(ctx, msgSend)
	if err != nil {
		return nil, err
	}

	return &v1.MsgWithdrawUnbondedResponse{}, nil
}

func executeMsg(ctx context.Context, msg proto.Message) ([]*codectypes.Any, error) {
	asAny, err := accountstd.PackAny(msg)
	if err != nil {
		return nil, err
	}

	return accountstd.ExecModuleAnys(ctx, []*codectypes.Any{asAny})
}
