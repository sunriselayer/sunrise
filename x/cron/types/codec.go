package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "sunrise/cron/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgAddSchedule{}, "sunrise/cron/MsgAddSchedule")
	legacy.RegisterAminoMsg(cdc, &MsgRemoveSchedule{}, "sunrise/cron/MsgRemoveSchedule")

	cdc.RegisterConcrete(&Params{}, "sunrise/cron/Params", nil)
}

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgAddSchedule{},
		&MsgRemoveSchedule{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
