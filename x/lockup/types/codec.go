package types

import (
	"cosmossdk.io/core/registry"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar registry.InterfaceRegistrar) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgInitLockupAccount{},
		&MsgNonVotingDelegate{},
		&MsgNonVotingUndelegate{},
		&MsgClaimRewards{},
		&MsgSend{},
	)

	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
