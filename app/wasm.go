package app

import (
	"path/filepath"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	"github.com/spf13/cast"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvm "github.com/CosmWasm/wasmvm/v2"
)

// registerWasmModules registers the Wasm modules for the application.
// https://github.com/CosmWasm/wasmd/blob/v0.60.0/app/app.go
// https://github.com/yerasyla/IgniteCLI-cosmwasm/blob/master/readme.md
func (app *App) registerWasmModules(appOpts servertypes.AppOptions, nodeConfig wasmtypes.NodeConfig) error {
	// set up non depinject support modules store keys
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(wasmtypes.StoreKey),
	); err != nil {
		return err
	}

	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	wasmDir := filepath.Join(homePath, "wasm")

	// https://ibc.cosmos.network/v8/ibc/light-clients/wasm/integration/
	// instantiate the Wasm VM with the chosen parameters
	wasmConfig := ibcwasmtypes.DefaultWasmConfig(DefaultNodeHome)
	wasmer, err := wasmvm.NewVM(
		wasmConfig.DataDir,
		wasmConfig.SupportedCapabilities,
		ibcwasmtypes.ContractMemoryLimit, // default of 32
		wasmConfig.ContractDebugMode,
		ibcwasmtypes.MemoryCacheSize,
	)
	if err != nil {
		return err
	}
	// create an Option slice (or append to an existing one)
	// with the option to use a custom Wasm VM instance
	wasmOpts := []wasmkeeper.Option{
		wasmkeeper.WithWasmEngine(wasmer),
	}
	app.WasmKeeper = wasmkeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(wasmtypes.StoreKey)),
		app.AuthKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		distrkeeper.NewQuerier(app.DistrKeeper),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		nodeConfig,
		wasmtypes.VMConfig{},
		wasmkeeper.BuiltInCapabilities(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)

	return nil
}

// Deprecated: Use BuiltInCapabilities from github.com/CosmWasm/wasmd/x/wasm/keeper
func AllCapabilities() []string {
	return wasmkeeper.BuiltInCapabilities()
}
