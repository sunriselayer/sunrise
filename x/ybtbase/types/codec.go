package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgCreate{},
		&MsgMint{},
		&MsgBurn{},
		&MsgAddYield{},
		&MsgGrantPermission{},
		&MsgRevokePermission{},
		&MsgClaimYield{},
		&MsgUpdateAdmin{},
		&MsgSend{},
	)

	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
