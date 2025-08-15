package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgPublishData{}, "sunrise/da/MsgPublishData")
	legacy.RegisterAminoMsg(cdc, &MsgSubmitInvalidity{}, "sunrise/da/MsgSubmitInvalidity")
	legacy.RegisterAminoMsg(cdc, &MsgSubmitValidityProof{}, "sunrise/da/MsgSubmitValidityProof")
	legacy.RegisterAminoMsg(cdc, &MsgRegisterProofDeputy{}, "sunrise/da/MsgRegisterProofDeputy")
	legacy.RegisterAminoMsg(cdc, &MsgUnregisterProofDeputy{}, "sunrise/da/MsgUnregisterProofDeputy")
	legacy.RegisterAminoMsg(cdc, &MsgVerifyData{}, "sunrise/da/MsgVerifyData")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "sunrise/da/MsgUpdateParams")

	cdc.RegisterConcrete(&Params{}, "sunrise/da/Params", nil)
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
