package app

import (
	"fmt"
	"io"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtos "github.com/cometbft/cometbft/libs/os"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	poolkeeper "github.com/cosmos/cosmos-sdk/x/protocolpool/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibcwasmkeeper "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/keeper"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/keeper"
	icahostkeeper "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	ibchookskeeper "github.com/sunriselayer/ibc-hooks/v10/keeper"

	"github.com/sunriselayer/sunrise/docs"
	damodulekeeper "github.com/sunriselayer/sunrise/x/da/keeper"
	feemodulekeeper "github.com/sunriselayer/sunrise/x/fee/keeper"
	liquidityincentivemodulekeeper "github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	liquiditypoolmodulekeeper "github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	lockupmodulekeeper "github.com/sunriselayer/sunrise/x/lockup/keeper"
	shareclassmodulekeeper "github.com/sunriselayer/sunrise/x/shareclass/keeper"
	stablemodulekeeper "github.com/sunriselayer/sunrise/x/stable/keeper"
	swapmodulekeeper "github.com/sunriselayer/sunrise/x/swap/keeper"
	swaptypes "github.com/sunriselayer/sunrise/x/swap/types"
	tokenconvertermodulekeeper "github.com/sunriselayer/sunrise/x/tokenconverter/keeper"

	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/sunriselayer/sunrise/app/gov"
	"github.com/sunriselayer/sunrise/app/mint"

	v1_0_1 "github.com/sunriselayer/sunrise/app/upgrades/v1_0_1"
)

const (
	AccountAddressPrefix = "sunrise"
	Name                 = "sunrise"
)

// DefaultNodeHome default home directories for the application daemon
var DefaultNodeHome string

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// keepers
	// only keepers required by the app are exposed
	// the list of all modules is available in the app_config
	AuthKeeper            authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	GovKeeper             *govkeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper
	CircuitBreakerKeeper  circuitkeeper.Keeper
	PoolKeeper            poolkeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper

	// ibc keepers
	IBCKeeper           *ibckeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper
	WasmClientKeeper    ibcwasmkeeper.Keeper

	// ibc Hooks
	IBCHooksKeeper ibchookskeeper.Keeper

	// cosmwasm keeper
	WasmKeeper wasmkeeper.Keeper

	DaKeeper                 damodulekeeper.Keeper
	FeeKeeper                feemodulekeeper.Keeper
	TokenconverterKeeper     tokenconvertermodulekeeper.Keeper
	ShareclassKeeper         shareclassmodulekeeper.Keeper
	LockupKeeper             lockupmodulekeeper.Keeper
	LiquiditypoolKeeper      liquiditypoolmodulekeeper.Keeper
	LiquidityincentiveKeeper liquidityincentivemodulekeeper.Keeper
	SwapKeeper               swapmodulekeeper.Keeper
	StableKeeper             stablemodulekeeper.Keeper
	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	var err error
	clienthelpers.EnvPrefix = Name
	DefaultNodeHome, err = clienthelpers.GetNodeHomeDirectory("." + Name)
	if err != nil {
		panic(err)
	}
}

// AppConfig returns the default app config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		appConfig,
		depinject.Supply(
			// supply custom module basics
			map[string]module.AppModuleBasic{
				genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			},
		),
		depinject.Provide(mint.ProvideMintFn),
		depinject.Provide(gov.ProvideCalculateVoteResultsAndVotingPowerFn),
	)
}

// New returns a reference to an initialized App.
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	var (
		app        = &App{}
		appBuilder *runtime.AppBuilder

		// merge the AppConfig and other configuration in one config
		appConfig = depinject.Configs(
			AppConfig(),
			depinject.Supply(
				appOpts, // supply app options
				logger,  // supply logger
				// here alternative options can be supplied to the DI container.
				// those options can be used f.e to override the default behavior of some modules.
				// for instance supplying a custom address codec for not using bech32 addresses.
				// read the depinject documentation and depinject module wiring for more information
				// on available options and how to use them.
			),
		)
	)

	var appModules map[string]appmodule.AppModule
	if err := depinject.Inject(appConfig,
		&appBuilder,
		&appModules,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AuthKeeper,
		&app.BankKeeper,
		&app.StakingKeeper,
		&app.SlashingKeeper,
		&app.MintKeeper,
		&app.DistrKeeper,
		&app.FeeGrantKeeper,
		&app.GovKeeper,
		&app.UpgradeKeeper,
		&app.AuthzKeeper,
		&app.ConsensusParamsKeeper,
		&app.CircuitBreakerKeeper,
		&app.PoolKeeper,
		&app.ParamsKeeper,
		&app.DaKeeper,
		&app.FeeKeeper,
		&app.TokenconverterKeeper,
		&app.ShareclassKeeper,
		&app.LockupKeeper,
		&app.LiquiditypoolKeeper,
		&app.LiquidityincentiveKeeper,
		&app.SwapKeeper,
		&app.StableKeeper,
	); err != nil {
		panic(err)
	}

	// add to default baseapp options
	// enable optimistic execution
	baseAppOptions = append(baseAppOptions, baseapp.SetOptimisticExecution())

	// build app
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// Register legacy modules
	wasmConfig, err := wasm.ReadNodeConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}
	if err := app.registerWasmAndIBCModules(appOpts, wasmConfig); err != nil {
		panic(err)
	}

	// register snapshot extensions
	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions(
			ibcwasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmClientKeeper),
		)
		if err != nil {
			panic(err)
		}
	}

	// <sunrise>
	app.SwapKeeper.TransferKeeper = &app.TransferKeeper
	// </sunrise>

	// register streaming services
	// if err := app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
	// panic(err)
	// }

	/****  Module Options ****/
	// <sunrise>
	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper: app.AuthKeeper,
				BankKeeper:    app.BankKeeper,
				ExtensionOptionChecker: func(a *codectypes.Any) bool {
					switch a.TypeUrl {
					case codectypes.MsgTypeURL(&swaptypes.SwapBeforeFeeExtension{}):
						return true
					default:
						return false
					}
				},
				SignModeHandler: app.txConfig.SignModeHandler(),
				FeegrantKeeper:  app.FeeGrantKeeper,
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper:     app.IBCKeeper,
			WasmConfig:    &wasmConfig,
			CircuitKeeper: &app.CircuitBreakerKeeper,
			FeeKeeper:     &app.FeeKeeper,
			SwapKeeper:    &app.SwapKeeper,
		},
	)
	if err != nil {
		panic(err)
	}
	app.SetAnteHandler(anteHandler)
	// </sunrise>

	// <sunrise>
	// Proposal handler for DA module
	daProposalHandler := NewProposalHandler(
		logger,
		app.DaKeeper,
		app.ModuleManager,
		baseapp.NewDefaultProposalHandler(app.Mempool(), app),
	)

	app.BaseApp.SetPrepareProposal(daProposalHandler.PrepareProposal())
	app.BaseApp.SetProcessProposal(daProposalHandler.ProcessProposal())
	app.BaseApp.SetPreBlocker(daProposalHandler.PreBlocker())
	// </sunrise>

	// create the simulation manager and define the order of the modules for deterministic simulations
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AuthKeeper, authsims.RandomGenesisAccounts, nil),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, overrideModules)

	app.sm.RegisterStoreDecoders()

	// A custom InitChainer sets if extra pre-init-genesis logic is required.
	// This is necessary for manually registered modules that do not support app wiring.
	// Manually set the module version map as shown below.
	// The upgrade module will automatically handle de-duplication of the module version map.
	app.SetInitChainer(func(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
		if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap()); err != nil {
			return nil, err
		}
		return app.App.InitChainer(ctx, req)
	})

	app.setupUpgradeHandlers()

	// <wasmd>
	// must be before Loading version
	// requires the snapshot store to be created and registered as a BaseAppOption
	// see cmd/wasmd/root.go: 206 - 214 approx
	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
	}
	// </wasmd>

	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	if loadLatest {
		ctx := app.BaseApp.NewUncachedContext(true, cmtproto.Header{})
		if err := app.WasmClientKeeper.InitializePinnedCodes(ctx); err != nil {
			cmtos.Exit(err.Error())
		}
		// <wasmd>
		// Initialize pinned codes in wasmvm as they are not persisted there
		if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			cmtos.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
		}
		// </wasmd>
	}

	return app
}

// GetSubspace returns a param subspace for a given module name.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// LegacyAmino returns App's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns App's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns App's InterfaceRegistry.
func (app *App) InterfaceRegistry() codectypes.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns App's TxConfig
func (app *App) TxConfig() client.TxConfig {
	return app.txConfig
}

// GetKey returns the KVStoreKey for the provided store key.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	kvStoreKey, ok := app.UnsafeFindStoreKey(storeKey).(*storetypes.KVStoreKey)
	if !ok {
		return nil
	}
	return kvStoreKey
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)
	// register swagger API in app.go so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}

	// register app's OpenAPI routes.
	docs.RegisterOpenAPIService(Name, apiSvr.Router)
}

// GetMaccPerms returns a copy of the module account permissions
//
// NOTE: This is solely to be used for testing purposes.
func GetMaccPerms() map[string][]string {
	dup := make(map[string][]string)
	for _, perms := range moduleAccPerms {
		dup[perms.Account] = perms.Permissions
	}

	return dup
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	result := make(map[string]bool)

	if len(blockAccAddrs) > 0 {
		for _, addr := range blockAccAddrs {
			result[addr] = true
		}
	} else {
		for addr := range GetMaccPerms() {
			result[addr] = true
		}
	}

	return result
}

func (app *App) setupUpgradeHandlers() {
	// Example upgrade handler.
	// When a planned upgrade height is reached, the old binary will panic and shut down, and the new binary
	// (which requires a PR) will take over with the new upgrade name defined below.
	// For more information, see: https://docs.cosmos.network/main/tooling/cosmovisor
	app.UpgradeKeeper.SetUpgradeHandler(
		v1_0_1.UpgradeName,
		v1_0_1.CreateUpgradeHandler(app.ModuleManager, app.Configurator()),
	)
}
