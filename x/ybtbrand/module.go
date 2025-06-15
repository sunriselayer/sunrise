package ybtbrand

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	ybtbrandtypes "github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the ybtbrand module.
type AppModuleBasic struct{}

// Name returns the ybtbrand module's name.
func (AppModuleBasic) Name() string {
	return ybtbrandtypes.ModuleName
}

// RegisterLegacyAminoCodec registers the ybtbrand module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// RegisterInterfaces registers the module's interface types
func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	ybtbrandtypes.RegisterInterfaces(registry)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the ybtbrand module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := ybtbrandtypes.RegisterQueryHandlerClient(clientCtx.CmdContext, mux, ybtbrandtypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// DefaultGenesis returns default genesis state as raw bytes for the ybtbrand module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(ybtbrandtypes.DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the ybtbrand module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState ybtbrandtypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ybtbrandtypes.ModuleName, err)
	}
	return genState.Validate()
}