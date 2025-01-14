package types

import (
	"cosmossdk.io/core/registry"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	v1 "github.com/sunriselayer/sunrise/x/accounts/self_delegatable_lockup/v1"
)

func RegisterInterfaces(registrar registry.InterfaceRegistrar) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgConvert{},
		&MsgSelfDelegate{},
		&MsgWithdrawSelfDelegationUnbonded{},
	)

	// TEMP: sunrise.accounts.self_delegatable_lockup.v1
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&v1.MsgSelfDelegate{},
		&v1.MsgWithdrawSelfDelegationUnbonded{},
		&v1.MsgSend{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
