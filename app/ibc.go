package app

import (
	"path/filepath"

	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cast"

	// IBC Wasm imports
	ibcwasm "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10"
	ibcwasmkeeper "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/keeper"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	icamodule "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/types"
	// ibctransfer "github.com/cosmos/ibc-go/v10/modules/apps/transfer"
	// ibctransferkeeper "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v10/modules/core"
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types" // nolint:staticcheck // Deprecated: params key table is needed for params migration
	ibcconnectiontypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	solomachine "github.com/cosmos/ibc-go/v10/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"

	"github.com/sunriselayer/sunrise/app/wasmclient"

	// transferv2 "github.com/cosmos/ibc-go/v10/modules/apps/transfer/v2"
	ibcapi "github.com/cosmos/ibc-go/v10/modules/core/api"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvm "github.com/CosmWasm/wasmvm/v2"

	// EVM imports
	srvflags "github.com/cosmos/evm/server/flags"

	ibctransfer "github.com/cosmos/evm/x/ibc/transfer" // NOTE: override ICS20 keeper to support IBC transfers of ERC20 tokens
	ibctransferkeeper "github.com/cosmos/evm/x/ibc/transfer/keeper"
	transferv2 "github.com/cosmos/evm/x/ibc/transfer/v2"

	"github.com/cosmos/evm/x/erc20"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	erc20v2 "github.com/cosmos/evm/x/erc20/v2"
	feemarket "github.com/cosmos/evm/x/feemarket"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	"github.com/cosmos/evm/x/precisebank"
	precisebankkeeper "github.com/cosmos/evm/x/precisebank/keeper"
	precisebanktypes "github.com/cosmos/evm/x/precisebank/types"
	evm "github.com/cosmos/evm/x/vm"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	swapmodule "github.com/sunriselayer/sunrise/x/swap/module"
	// this line is used by starport scaffolding # ibc/app/import
)

// registerWasmAndIBCModules register CosmWasm and IBC keepers and non dependency inject modules.
func (app *App) registerWasmAndIBCModules(appOpts servertypes.AppOptions, nodeConfig wasmtypes.NodeConfig) error {
	// set up non depinject support modules store keys
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(ibcexported.StoreKey),
		storetypes.NewKVStoreKey(ibctransfertypes.StoreKey),
		storetypes.NewKVStoreKey(icahosttypes.StoreKey),
		storetypes.NewKVStoreKey(icacontrollertypes.StoreKey),
		storetypes.NewKVStoreKey(ibcwasmtypes.StoreKey),
		storetypes.NewKVStoreKey(wasmtypes.StoreKey),
		storetypes.NewKVStoreKey(evmtypes.StoreKey),
		storetypes.NewKVStoreKey(feemarkettypes.StoreKey),
		storetypes.NewKVStoreKey(erc20types.StoreKey),
		storetypes.NewKVStoreKey(precisebanktypes.StoreKey),
	); err != nil {
		return err
	}

	// register the key tables for legacy param subspaces
	keyTable := ibcclienttypes.ParamKeyTable()
	keyTable.RegisterParamSet(&ibcconnectiontypes.Params{})
	app.ParamsKeeper.Subspace(ibcexported.ModuleName).WithKeyTable(keyTable)
	app.ParamsKeeper.Subspace(ibctransfertypes.ModuleName).WithKeyTable(ibctransfertypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(icacontrollertypes.SubModuleName).WithKeyTable(icacontrollertypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(icahosttypes.SubModuleName).WithKeyTable(icahosttypes.ParamKeyTable())

	govModuleAddr, _ := app.AuthKeeper.AddressCodec().BytesToString(authtypes.NewModuleAddress(govtypes.ModuleName))

	// Create IBC keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(ibcexported.StoreKey)),
		app.GetSubspace(ibcexported.ModuleName),
		app.UpgradeKeeper,
		govModuleAddr,
	)

	// Create IBC transfer keeper
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(ibctransfertypes.StoreKey)),
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.MsgServiceRouter(),
		app.AuthKeeper,
		app.BankKeeper,
		govModuleAddr,
	)

	// Create interchain account keepers
	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(icahosttypes.StoreKey)),
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper, // ICS4Wrapper
		app.IBCKeeper.ChannelKeeper,
		app.AuthKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		govModuleAddr,
	)

	app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(icacontrollertypes.StoreKey)),
		app.GetSubspace(icacontrollertypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.MsgServiceRouter(),
		govModuleAddr,
	)

	// create IBC module from bottom to top of stack
	var (
		transferStack      porttypes.IBCModule = ibctransfer.NewIBCModule(app.TransferKeeper)
		transferStackV2    ibcapi.IBCModule    = transferv2.NewIBCModule(app.TransferKeeper)
		icaControllerStack porttypes.IBCModule = icacontroller.NewIBCMiddleware(app.ICAControllerKeeper)
		icaHostStack       porttypes.IBCModule = icahost.NewIBCModule(app.ICAHostKeeper)
	)

	// <sunrise>
	transferStack = swapmodule.NewIBCMiddleware(transferStack, &app.SwapKeeper)
	// </sunrise>

	// <wasmd>
	// https://github.com/CosmWasm/wasmd/blob/v0.60.0/app/app.go
	// https://github.com/yerasyla/IgniteCLI-cosmwasm/blob/master/readme.md
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	wasmDir := filepath.Join(homePath, "wasm")
	// https://ibc.cosmos.network/v8/ibc/light-clients/wasm/integration/
	// instantiate the Wasm VM with the chosen parameters
	wasmConfig := ibcwasmtypes.DefaultWasmConfig(DefaultNodeHome)
	wasmer, err := wasmvm.NewVM(
		wasmConfig.DataDir,
		wasmkeeper.BuiltInCapabilities(), //  wasmConfig.SupportedCapabilities support only `iterator`
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

	// Create fee enabled wasm ibc Stack
	wasmStack := wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper)
	// </wasmd>

	// <evmd>
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		app.appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		storetypes.NewTransientStoreKey(feemarkettypes.TransientKey),
	)

	app.PreciseBankKeeper = precisebankkeeper.NewKeeper(
		appCodec,
		keys[precisebanktypes.StoreKey],
		app.BankKeeper,
		app.AccountKeeper,
	)

	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))
	app.EvmKeeper = evmkeeper.NewKeeper(
		app.appCodec,
		keys,
		tkeys,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper, // StakingKeeper can be required for certain precompiles
		app.FeeMarketKeeper,
		tracer,
	)

	app.Erc20Keeper = erc20keeper.NewKeeper(
		keys[erc20types.StoreKey],
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper,
		app.PreciseBankKeeper,
		app.EVMKeeper,
		app.StakingKeeper,
		&app.TransferKeeper,
	)

	// NOTE: we are adding all available Cosmos EVM EVM extensions.
	// Not all of them need to be enabled, which can be configured on a per-chain basis.
	app.EvmKeeper.WithStaticPrecompiles()

	transferStack = erc20.NewIBCMiddleware(app.Erc20Keeper, transferStack)
	transferStackV2 = erc20v2.NewIBCMiddleware(transferStackV2, app.Erc20Keeper)
	// </evmd>

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		// <wasmd>
		AddRoute(wasmtypes.ModuleName, wasmStack).
		// </wasmd>
		AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
		AddRoute(icahosttypes.SubModuleName, icaHostStack)

	// this line is used by starport scaffolding # ibc/app/module

	// Seal the IBC Router
	app.IBCKeeper.SetRouter(ibcRouter)

	// <sunrise>
	ibcRouterV2 := ibcapi.NewRouter().
		AddRoute(ibctransfertypes.PortID, transferStackV2)
	app.IBCKeeper.SetRouterV2(ibcRouterV2)
	// </sunrise>

	storeProvider := app.IBCKeeper.ClientKeeper.GetStoreProvider()
	tmLightClientModule := ibctm.NewLightClientModule(app.appCodec, storeProvider)
	soloLightClientModule := solomachine.NewLightClientModule(app.appCodec, storeProvider)

	// <sunrise>
	wasmLightClientQuerier := ibcwasmkeeper.QueryPlugins{
		Stargate: ibcwasmkeeper.AcceptListStargateQuerier([]string{
			"/ibc.core.client.v1.Query/ClientState",
			"/ibc.core.client.v1.Query/ConsensusState",
			"/ibc.core.connection.v1.Query/Connection",
		}, app.GRPCQueryRouter()),
		Custom: wasmclient.CustomQuerier(),
	}

	app.WasmClientKeeper = ibcwasmkeeper.NewKeeperWithVM(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(ibcwasmtypes.StoreKey)),
		app.IBCKeeper.ClientKeeper,
		govModuleAddr,
		wasmer,
		app.GRPCQueryRouter(),
		ibcwasmkeeper.WithQueryPlugins(&wasmLightClientQuerier),
	)
	// </sunrise>

	wasmLightClientModule := ibcwasm.NewLightClientModule(app.WasmClientKeeper, storeProvider)
	app.IBCKeeper.ClientKeeper.AddRoute(ibctm.ModuleName, &tmLightClientModule)
	app.IBCKeeper.ClientKeeper.AddRoute(solomachine.ModuleName, &soloLightClientModule)
	app.IBCKeeper.ClientKeeper.AddRoute(ibcwasmtypes.ModuleName, &wasmLightClientModule)

	// register IBC modules
	if err := app.RegisterModules(
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.TransferKeeper),
		icamodule.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		ibctm.NewAppModule(tmLightClientModule),
		solomachine.NewAppModule(soloLightClientModule),
		ibcwasm.NewAppModule(app.WasmClientKeeper),
		wasm.NewAppModule(app.appCodec, &app.WasmKeeper, app.StakingKeeper, app.AuthKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		evm.NewAppModule(app.EvmKeeper, app.AuthKeeper, app.BankKeeper),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		erc20.NewAppModule(app.Erc20Keeper, app.EvmKeeper, app.PreciseBankKeeper, app.TransferKeeper),
		precisebank.NewAppModule(app.PreciseBankKeeper, app.EvmKeeper, app.BankKeeper),
	); err != nil {
		return err
	}

	return nil
}

// Since the IBC modules don't support dependency injection, we need to
// manually register the modules on the client side.
// This needs to be removed after IBC supports App Wiring.
func RegisterWasmAndIBC(cdc codec.Codec, registry cdctypes.InterfaceRegistry) map[string]appmodule.AppModule {
	modules := map[string]appmodule.AppModule{
		ibcexported.ModuleName:      ibc.NewAppModule(&ibckeeper.Keeper{}),
		ibctransfertypes.ModuleName: ibctransfer.NewAppModule(ibctransferkeeper.Keeper{}),
		icatypes.ModuleName:         icamodule.NewAppModule(&icacontrollerkeeper.Keeper{}, &icahostkeeper.Keeper{}),
		ibctm.ModuleName:            ibctm.NewAppModule(ibctm.NewLightClientModule(cdc, ibcclienttypes.StoreProvider{})),
		solomachine.ModuleName:      solomachine.NewAppModule(solomachine.NewLightClientModule(cdc, ibcclienttypes.StoreProvider{})),
		ibcwasmtypes.ModuleName:     ibcwasm.NewAppModule(ibcwasmkeeper.Keeper{}),
		wasmtypes.ModuleName:        wasm.NewAppModule(cdc, &wasmkeeper.Keeper{}, nil, nil, nil, nil, nil),
		evmtypes.ModuleName:         evm.NewAppModule(&evmkeeper.Keeper{}, nil, nil),
		feemarkettypes.ModuleName:   feemarket.NewAppModule(feemarketkeeper.Keeper{}),
		erc20types.ModuleName:       erc20.NewAppModule(erc20keeper.Keeper{}, &evmkeeper.Keeper{}, precisebankkeeper.Keeper{}, ibctransferkeeper.Keeper{}),
		precisebanktypes.ModuleName: precisebank.NewAppModule(precisebankkeeper.Keeper{}, &evmkeeper.Keeper{}, nil),
	}

	for _, m := range modules {
		if mr, ok := m.(module.AppModuleBasic); ok {
			mr.RegisterInterfaces(registry)
		}
	}

	return modules
}
