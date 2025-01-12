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
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	selfdelegationproxy "github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy"
	selfdelegationproxytypes "github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy/v1"
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

	// If the amount is in the fee denom
	if msg.Amount.Denom == params.FeeDenom {
		delegator, err := m.addressCodec.StringToBytes(msg.DelegatorAddress)
		if err != nil {
			return nil, err
		}
		accType, err := m.accountKeeper.AccountsByType.Get(ctx, delegator)
		if err != nil {
			return nil, err
		}

		var rootOwner []byte
		switch accType {
		// Case of lockup accounts
		case lockup.CONTINUOUS_LOCKING_ACCOUNT,
			lockup.DELAYED_LOCKING_ACCOUNT,
			lockup.PERIODIC_LOCKING_ACCOUNT,
			lockup.PERMANENT_LOCKING_ACCOUNT:

			res, err := m.accountKeeper.Query(ctx, delegator, &lockuptypes.QueryLockupAccountInfoRequest{})
			if err != nil {
				return nil, err
			}
			rootOwner, err = m.addressCodec.StringToBytes(res.(*lockuptypes.QueryLockupAccountInfoResponse).Owner)
			if err != nil {
				return nil, err
			}
		default:
			rootOwner = delegator
		}
		// Convert to val address
		rootOwnerValAddress, err := m.valAddressCodec.BytesToString(rootOwner)
		if err != nil {
			return nil, err
		}

		if rootOwnerValAddress != msg.ValidatorAddress {
			return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "delegation with denom %s is only for self delegation", params.FeeDenom)
		}

		// Check proxy account existence
		has, err := m.tokenconverterKeeper.SelfDelegationProxy.Has(ctx, delegator)
		if err != nil {
			return nil, err
		}
		var proxyAddrBytes []byte
		if has {
			proxyAddrBytes, err = m.tokenconverterKeeper.SelfDelegationProxy.Get(ctx, delegator)
			if err != nil {
				return nil, err
			}
		} else {
			rootOwnerString, err := m.addressCodec.BytesToString(rootOwner)
			if err != nil {
				return nil, err
			}

			// Create proxy account
			_, proxyAddrBytes, err = m.accountKeeper.Init(
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
			err = m.tokenconverterKeeper.SelfDelegationProxy.Set(ctx, delegator, proxyAddrBytes)
			if err != nil {
				return nil, err
			}
		}

		// ConvertReverse
		err = m.tokenconverterKeeper.Convert(ctx, msg.Amount.Amount, proxyAddrBytes)
		if err != nil {
			return nil, err
		}

		// Delegate from proxy account
		proxyAddr, err := m.addressCodec.BytesToString(proxyAddrBytes)
		if err != nil {
			return nil, err
		}
		res, err := m.Keeper.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgDelegate{
			DelegatorAddress: proxyAddr,
			ValidatorAddress: rootOwnerValAddress,
			Amount:           sdk.NewCoin(params.BondDenom, msg.Amount.Amount),
		})

		return res.(*stakingtypes.MsgDelegateResponse), err
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
