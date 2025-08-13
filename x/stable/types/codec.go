package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "sunrise/stable/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgMint{}, "sunrise/MsgMint")
	legacy.RegisterAminoMsg(cdc, &MsgBurn{}, "sunrise/MsgBurn")
}

// RegisterInterfaces registers the module's interface types
func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBurn{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMint{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
