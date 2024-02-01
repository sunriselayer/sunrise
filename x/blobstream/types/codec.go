package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

const URLMsgRegisterEVMAddress = "/sunrise.blob.v1.MsgRegisterEVMAddress"

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterEvmAddress{}, URLMsgRegisterEVMAddress, nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterEvmAddress{},
		&MsgUpdateParams{},
	)
	// this line is used by starport scaffolding # 3

	registry.RegisterInterface(
		"AttestationRequestI",
		(*AttestationRequestI)(nil),
		&DataCommitment{},
		&Valset{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
