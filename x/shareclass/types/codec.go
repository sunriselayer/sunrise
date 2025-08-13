package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "sunrise/shareclass/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgNonVotingDelegate{}, "sunrise/MsgNonVotingDelegate")
	legacy.RegisterAminoMsg(cdc, &MsgNonVotingUndelegate{}, "sunrise/MsgNonVotingUndelegate")
	legacy.RegisterAminoMsg(cdc, &MsgClaimRewards{}, "sunrise/MsgClaimRewards")
	legacy.RegisterAminoMsg(cdc, &MsgCreateValidator{}, "sunrise/MsgCreateValidator")
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
