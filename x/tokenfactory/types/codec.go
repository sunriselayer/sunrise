package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateDenom{}, "sunrise/tokenfactory/MsgCreateDenom")
	legacy.RegisterAminoMsg(cdc, &MsgMint{}, "sunrise/tokenfactory/MsgMint")
	legacy.RegisterAminoMsg(cdc, &MsgBurn{}, "sunrise/tokenfactory/MsgBurn")
	legacy.RegisterAminoMsg(cdc, &MsgChangeAdmin{}, "sunrise/tokenfactory/MsgChangeAdmin")
	legacy.RegisterAminoMsg(cdc, &MsgSetDenomMetadata{}, "sunrise/tokenfactory/MsgSetDenomMetadata")
	legacy.RegisterAminoMsg(cdc, &MsgSetBeforeSendHook{}, "sunrise/tokenfactory/MsgSetBeforeSendHook")
	legacy.RegisterAminoMsg(cdc, &MsgForceTransfer{}, "sunrise/tokenfactory/MsgForceTransfer")
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
