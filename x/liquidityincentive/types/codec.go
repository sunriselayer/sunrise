package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "sunrise/incentive/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgStartNewEpoch{}, "sunrise/incentive/MsgStartNewEpoch")
	legacy.RegisterAminoMsg(cdc, &MsgVoteGauge{}, "sunrise/incentive/MsgVoteGauge")
	legacy.RegisterAminoMsg(cdc, &MsgRegisterBribe{}, "sunrise/incentive/MsgRegisterBribe")
	legacy.RegisterAminoMsg(cdc, &MsgClaimBribes{}, "sunrise/incentive/MsgClaimBribes")

	cdc.RegisterConcrete(&Params{}, "sunrise/incentive/Params", nil)
	cdc.RegisterConcrete(&Gauge{}, "sunrise/incentive/Gauge", nil)
	cdc.RegisterConcrete(&Bribe{}, "sunrise/incentive/Bribe", nil)
}

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVoteGauge{},
		&MsgRegisterBribe{},
		&MsgClaimBribes{},
		&MsgStartNewEpoch{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
