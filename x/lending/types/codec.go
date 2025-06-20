package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgLiquidate{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRepay{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBorrow{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSupply{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
