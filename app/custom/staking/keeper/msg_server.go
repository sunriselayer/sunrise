package keeper

import (
	"context"

	stakingkeeper "cosmossdk.io/x/staking/keeper"
	stakingtypes "cosmossdk.io/x/staking/types"

	customtypes "github.com/sunriselayer/sunrise/app/custom/types"
)

type msgServer struct {
	stakingKeeper        customtypes.StakingKeeper
	tokenconverterKeeper customtypes.TokenConverterKeeper
}

func NewMsgServerImpl(
	stakingKeeper customtypes.StakingKeeper,
	tokenconverterKeeper customtypes.TokenConverterKeeper,
) stakingtypes.MsgServer {
	return &msgServer{
		stakingKeeper:        stakingKeeper,
		tokenconverterKeeper: tokenconverterKeeper,
	}
}

func (m msgServer) StakingMsgServer() stakingtypes.MsgServer {
	return stakingkeeper.NewMsgServerImpl(m.stakingKeeper.(*stakingkeeper.Keeper))
}

func (m msgServer) CreateValidator(ctx context.Context, msg *stakingtypes.MsgCreateValidator) (*stakingtypes.MsgCreateValidatorResponse, error) {
	return m.StakingMsgServer().CreateValidator(ctx, msg)
}

func (m msgServer) EditValidator(ctx context.Context, msg *stakingtypes.MsgEditValidator) (*stakingtypes.MsgEditValidatorResponse, error) {
	return m.StakingMsgServer().EditValidator(ctx, msg)
}

func (m msgServer) Delegate(ctx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error) {
	feeDenom, err := m.tokenconverterKeeper.GetFeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	// If the amount is in the fee denom
	if msg.Amount.Denom == feeDenom {
		return m.tokenconverterKeeper.SelfDelegate(ctx, msg)
	}
	return m.StakingMsgServer().Delegate(ctx, msg)
}

func (m msgServer) BeginRedelegate(ctx context.Context, msg *stakingtypes.MsgBeginRedelegate) (*stakingtypes.MsgBeginRedelegateResponse, error) {
	return m.StakingMsgServer().BeginRedelegate(ctx, msg)
}

func (m msgServer) Undelegate(ctx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error) {
	return m.StakingMsgServer().Undelegate(ctx, msg)
}

func (m msgServer) CancelUnbondingDelegation(ctx context.Context, msg *stakingtypes.MsgCancelUnbondingDelegation) (*stakingtypes.MsgCancelUnbondingDelegationResponse, error) {
	return m.StakingMsgServer().CancelUnbondingDelegation(ctx, msg)
}

func (m msgServer) UpdateParams(ctx context.Context, msg *stakingtypes.MsgUpdateParams) (*stakingtypes.MsgUpdateParamsResponse, error) {
	return m.StakingMsgServer().UpdateParams(ctx, msg)
}

func (m msgServer) RotateConsPubKey(ctx context.Context, msg *stakingtypes.MsgRotateConsPubKey) (*stakingtypes.MsgRotateConsPubKeyResponse, error) {
	return m.StakingMsgServer().RotateConsPubKey(ctx, msg)
}
