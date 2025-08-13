package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "sunrise/lockup/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgInitLockupAccount{}, "sunrise/MsgInitLockupAccount")
	legacy.RegisterAminoMsg(cdc, &MsgNonVotingDelegate{}, "sunrise/lockup/MsgNonVotingDelegate")
	legacy.RegisterAminoMsg(cdc, &MsgNonVotingUndelegate{}, "sunrise/lockup/MsgNonVotingUndelegate")
	legacy.RegisterAminoMsg(cdc, &MsgClaimRewards{}, "sunrise/MsgClaimRewards")
	legacy.RegisterAminoMsg(cdc, &MsgSend{}, "sunrise/lockup/MsgSend")
}

// RegisterInterfaces registers the module's interface types
func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgInitLockupAccount{},
		&MsgNonVotingDelegate{},
		&MsgNonVotingUndelegate{},
		&MsgClaimRewards{},
		&MsgSend{},
	)

	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
