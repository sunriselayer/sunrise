package keeper

import (
	"context"

	"cosmossdk.io/core/address"
	errorsmod "cosmossdk.io/errors"
	accounts "cosmossdk.io/x/accounts"
	"cosmossdk.io/x/accounts/defaults/lockup"
	lockuptypes "cosmossdk.io/x/accounts/defaults/lockup/v1"
	stakingkeeper "cosmossdk.io/x/staking/keeper"
	stakingtypes "cosmossdk.io/x/staking/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	tokenconverterkeeper "github.com/sunriselayer/sunrise/x/tokenconverter/keeper"
)

type msgServer struct {
	*stakingkeeper.Keeper
	addressCodec         address.Codec
	valAddressCodec      address.ValidatorAddressCodec
	accountKeeper        accounts.Keeper
	tokenconverterKeeper tokenconverterkeeper.Keeper
}

func NewMsgServerImpl(keeper *stakingkeeper.Keeper, tokenconverterKeeper tokenconverterkeeper.Keeper) stakingtypes.MsgServer {
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

	if msg.Amount.Denom == params.FeeDenom {
		delegator, err := m.addressCodec.StringToBytes(msg.DelegatorAddress)
		if err != nil {
			return nil, err
		}
		accType, err := m.accountKeeper.AccountsByType.Get(ctx, delegator)
		if err != nil {
			return nil, err
		}

		var owner []byte
		switch accType {
		case lockup.CONTINUOUS_LOCKING_ACCOUNT,
			lockup.DELAYED_LOCKING_ACCOUNT,
			lockup.PERIODIC_LOCKING_ACCOUNT,
			lockup.PERMANENT_LOCKING_ACCOUNT:

			resI, err := m.accountKeeper.Query(ctx, delegator, &lockuptypes.QueryLockupAccountInfoRequest{})
			if err != nil {
				return nil, err
			}
			res, _ := resI.(*lockuptypes.QueryLockupAccountInfoResponse)
			owner, err = m.addressCodec.StringToBytes(res.Owner)
			if err != nil {
				return nil, err
			}
		default:
			owner = delegator
		}
		ownerValAddress, err := m.valAddressCodec.BytesToString(owner)
		if err != nil {
			return nil, err
		}

		if ownerValAddress != msg.ValidatorAddress {
			return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "delegation with denom %s is only for self delegation", params.FeeDenom)
		}

		// Check proxy account existence
		var exist bool
		if !exist {
			// Create proxy account

		}

		// ConvertReverse

		// Move tokens to proxy account

		// Delegate from proxy account

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
