package keeper

import (
	"context"

	stakingkeeper "cosmossdk.io/x/staking/keeper"
	stakingtypes "cosmossdk.io/x/staking/types"

	tokenconverterkeeper "github.com/sunriselayer/sunrise/x/tokenconverter/keeper"
)

type msgServer struct {
	*stakingkeeper.Keeper
	tokenconverterKeeper *tokenconverterkeeper.Keeper
}

func NewMsgServerImpl(keeper *stakingkeeper.Keeper, tokenconverterKeeper *tokenconverterkeeper.Keeper) stakingtypes.MsgServer {
	return &msgServer{
		Keeper:               keeper,
		tokenconverterKeeper: tokenconverterKeeper,
	}
}

func (m msgServer) CreateValidator(ctx context.Context, msg *stakingtypes.MsgCreateValidator) (*stakingtypes.MsgCreateValidatorResponse, error) {
	return stakingkeeper.NewMsgServerImpl(m.Keeper).CreateValidator(ctx, msg)
}

func (m msgServer) EditValidator(ctx context.Context, msg *stakingtypes.MsgEditValidator) (*stakingtypes.MsgEditValidatorResponse, error) {
	return stakingkeeper.NewMsgServerImpl(m.Keeper).EditValidator(ctx, msg)
}

func (m msgServer) Delegate(ctx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error) {
	params, err := m.tokenconverterKeeper.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// If the amount is in the fee denom
	if msg.Amount.Denom == params.FeeDenom {
		return m.tokenconverterKeeper.SelfDelegate(ctx, msg)
	}
	return stakingkeeper.NewMsgServerImpl(m.Keeper).Delegate(ctx, msg)
}

func (m msgServer) BeginRedelegate(ctx context.Context, msg *stakingtypes.MsgBeginRedelegate) (*stakingtypes.MsgBeginRedelegateResponse, error) {
	return stakingkeeper.NewMsgServerImpl(m.Keeper).BeginRedelegate(ctx, msg)
}

func (m msgServer) Undelegate(ctx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error) {
	return stakingkeeper.NewMsgServerImpl(m.Keeper).Undelegate(ctx, msg)
}

func (m msgServer) CancelUnbondingDelegation(ctx context.Context, msg *stakingtypes.MsgCancelUnbondingDelegation) (*stakingtypes.MsgCancelUnbondingDelegationResponse, error) {
	return stakingkeeper.NewMsgServerImpl(m.Keeper).CancelUnbondingDelegation(ctx, msg)
}

func (m msgServer) UpdateParams(ctx context.Context, msg *stakingtypes.MsgUpdateParams) (*stakingtypes.MsgUpdateParamsResponse, error) {
	return stakingkeeper.NewMsgServerImpl(m.Keeper).UpdateParams(ctx, msg)
}

func (m msgServer) RotateConsPubKey(ctx context.Context, msg *stakingtypes.MsgRotateConsPubKey) (*stakingtypes.MsgRotateConsPubKeyResponse, error) {
	return stakingkeeper.NewMsgServerImpl(m.Keeper).RotateConsPubKey(ctx, msg)
}
