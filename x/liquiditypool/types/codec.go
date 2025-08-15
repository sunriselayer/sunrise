package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "sunrise/pool/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgCreatePool{}, "sunrise/MsgCreatePool")
	legacy.RegisterAminoMsg(cdc, &MsgCreatePosition{}, "sunrise/MsgCreatePosition")
	legacy.RegisterAminoMsg(cdc, &MsgIncreaseLiquidity{}, "sunrise/MsgIncreaseLiquidity")
	legacy.RegisterAminoMsg(cdc, &MsgDecreaseLiquidity{}, "sunrise/MsgDecreaseLiquidity")
	legacy.RegisterAminoMsg(cdc, &MsgClaimRewards{}, "sunrise/pool/MsgClaimRewards")

	cdc.RegisterConcrete(&Params{}, "sunrise/pool/Params", nil)
	cdc.RegisterConcrete(&Pool{}, "sunrise/pool/Pool", nil)
}

// RegisterInterfaces registers the module's interface types
func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePool{},
		&MsgCreatePosition{},
		&MsgIncreaseLiquidity{},
		&MsgDecreaseLiquidity{},
		&MsgClaimRewards{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
