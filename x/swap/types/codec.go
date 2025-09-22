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
	legacy.RegisterAminoMsg(cdc, &MsgSwapExactAmountIn{}, "sunrise/swap/MsgSwapExactAmountIn")
	legacy.RegisterAminoMsg(cdc, &MsgSwapExactAmountOut{}, "sunrise/swap/MsgSwapExactAmountOut")

	// register Route's oneof field's interface
	cdc.RegisterConcrete(&Params{}, "sunrise/swap/Params", nil)
	cdc.RegisterInterface((*isRoute_Strategy)(nil), nil)
	cdc.RegisterConcrete(&Route_Pool{}, "sunrise/swap/RoutePool", nil)
	cdc.RegisterConcrete(&Route_Series{}, "sunrise/swap/RouteSeries", nil)
	cdc.RegisterConcrete(&Route_Parallel{}, "sunrise/swap/RouteParallel", nil)
	cdc.RegisterConcrete(&Route{}, "sunrise/swap/Route", nil)
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
