package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgPublishData{}, "sunrise/MsgPublishData")
	legacy.RegisterAminoMsg(cdc, &MsgSubmitInvalidity{}, "sunrise/MsgSubmitInvalidity")
	legacy.RegisterAminoMsg(cdc, &MsgSubmitValidityProof{}, "sunrise/MsgSubmitValidityProof")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "sunrise/da/MsgUpdateParams")
}

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPublishData{},
		&MsgSubmitInvalidity{},
		&MsgSubmitValidityProof{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
