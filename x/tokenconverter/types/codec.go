package types

import (
	"cosmossdk.io/core/registry"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	lockuptypes "github.com/sunriselayer/sunrise/x/accounts/self_delegatable_lockup/v1"
	proxytypes "github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy/v1"
)

func RegisterInterfaces(registrar registry.InterfaceRegistrar) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgConvert{},
		&MsgSelfDelegate{},
		&MsgWithdrawSelfDelegationUnbonded{},
	)

	// TEMP: sunrise.accounts.self_delegatable_lockup.v1
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&lockuptypes.MsgInitSelfDelegatableLockupAccount{},
		&lockuptypes.MsgInitSelfDelegatableLockupAccountResponse{},
		&lockuptypes.MsgSelfDelegate{},
		&lockuptypes.MsgWithdrawSelfDelegationUnbonded{},
		&lockuptypes.MsgSend{},
	)

	// TEMP: sunrise.accounts.self_delegation_proxy.v1
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&proxytypes.MsgInit{},
		&proxytypes.MsgUndelegate{},
		&proxytypes.MsgWithdrawReward{},
		&proxytypes.MsgSend{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
