package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "sunrise/swap/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgSwapExactAmountIn{}, "sunrise/MsgSwapExactAmountIn")
	legacy.RegisterAminoMsg(cdc, &MsgSwapExactAmountOut{}, "sunrise/MsgSwapExactAmountOut")

	cdc.RegisterConcrete(&Params{}, "sunrise/swap/Params", nil)
	cdc.RegisterConcrete(&SwapBeforeFeeExtension{}, "sunrise/swap/SwapBeforeFeeExtension", nil)
}

// RegisterInterfaces registers the module's interface types
func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSwapExactAmountIn{},
		&MsgSwapExactAmountOut{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)

	// Register extension options
	registrar.RegisterImplementations(
		(*tx.TxExtensionOptionI)(nil),
		&SwapBeforeFeeExtension{},
	)

	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
