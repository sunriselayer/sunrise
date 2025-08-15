package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "sunrise/sc/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgNonVotingDelegate{}, "sunrise/sc/MsgNonVotingDelegate")
	legacy.RegisterAminoMsg(cdc, &MsgNonVotingUndelegate{}, "sunrise/sc/MsgNonVotingUndelegate")
	legacy.RegisterAminoMsg(cdc, &MsgClaimRewards{}, "sunrise/sc/MsgClaimRewards")
	legacy.RegisterAminoMsg(cdc, &MsgCreateValidator{}, "sunrise/sc/MsgCreateValidator")

	cdc.RegisterConcrete(&Params{}, "sunrise/sc/Params", nil)
}

// RegisterInterfaces registers the module's interface types
func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateValidator{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNonVotingDelegate{},
		&MsgNonVotingUndelegate{},
		&MsgClaimRewards{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
