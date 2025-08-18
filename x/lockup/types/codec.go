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
	legacy.RegisterAminoMsg(cdc, &MsgInitLockupAccount{}, "sunrise/lockup/MsgInitLockupAccount")
	legacy.RegisterAminoMsg(cdc, &MsgNonVotingDelegate{}, "sunrise/lockup/MsgNonVotingDelegate")
	legacy.RegisterAminoMsg(cdc, &MsgNonVotingUndelegate{}, "sunrise/lockup/MsgNonVotingUndelegate")
	legacy.RegisterAminoMsg(cdc, &MsgClaimRewards{}, "sunrise/lockup/MsgClaimRewards")
	legacy.RegisterAminoMsg(cdc, &MsgSend{}, "sunrise/lockup/MsgSend")

	cdc.RegisterConcrete(&Params{}, "sunrise/lockup/Params", nil)
	cdc.RegisterConcrete(&LockupAccount{}, "sunrise/lockup/LockupAccount", nil)
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
