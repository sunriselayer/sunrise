package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateDenom{}, "sunrise/factory/MsgCreateDenom")
	legacy.RegisterAminoMsg(cdc, &MsgMint{}, "sunrise/factory/MsgMint")
	legacy.RegisterAminoMsg(cdc, &MsgBurn{}, "sunrise/factory/MsgBurn")
	legacy.RegisterAminoMsg(cdc, &MsgChangeAdmin{}, "sunrise/factory/MsgChangeAdmin")
	legacy.RegisterAminoMsg(cdc, &MsgSetDenomMetadata{}, "sunrise/factory/MsgSetDenomMetadata")
	legacy.RegisterAminoMsg(cdc, &MsgSetBeforeSendHook{}, "sunrise/factory/MsgSetBeforeSendHook")
	legacy.RegisterAminoMsg(cdc, &MsgForceTransfer{}, "sunrise/factory/MsgForceTransfer")
}

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateDenom{},
		&MsgMint{},
		&MsgBurn{},
		&MsgChangeAdmin{},
		&MsgSetDenomMetadata{},
		&MsgSetBeforeSendHook{},
		&MsgForceTransfer{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
