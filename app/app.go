package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
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
	testdata_pulsar "github.com/cosmos/cosmos-sdk/testutil/testdata/testpb"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"github.com/skip-mev/block-sdk/v2/abci"
	"github.com/skip-mev/block-sdk/v2/abci/checktx"
	"github.com/skip-mev/block-sdk/v2/block"
	"github.com/skip-mev/block-sdk/v2/block/base"
	"github.com/skip-mev/block-sdk/v2/block/service"
	mevlane "github.com/skip-mev/block-sdk/v2/lanes/mev"
	"github.com/sunriselayer/sunrise/app/ante"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	defaultoverrides "github.com/sunriselayer/sunrise/app/defaultoverrides"
	feetypes "github.com/sunriselayer/sunrise/x/fee/types"
	tokenconvertertypes "github.com/sunriselayer/sunrise/x/tokenconverter/types"

	// this line is used by starport scaffolding # stargate/app/moduleImport

	"github.com/sunriselayer/sunrise/app/keepers"
	"github.com/sunriselayer/sunrise/app/upgrades"
	v0_1_5_test "github.com/sunriselayer/sunrise/app/upgrades/v0.1.5-test"
	"github.com/sunriselayer/sunrise/docs"
)

const (
	Bech32MainPrefix = "sunrise"
	Name             = "sunrise"

	Bech32PrefixAccAddr  = Bech32MainPrefix
	Bech32PrefixAccPub   = Bech32MainPrefix + "pub"
	Bech32PrefixValAddr  = Bech32MainPrefix + "valoper"
	Bech32PrefixValPub   = Bech32MainPrefix + "valoperpub"
	Bech32PrefixConsAddr = Bech32MainPrefix + "valcons"
	Bech32PrefixConsPub  = Bech32MainPrefix + "valconspub"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// <sunrise>
	Upgrades = []upgrades.Upgrade{v0_1_5_test.Upgrade}
	// </sunrise>
)

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*runtime.App
	// <sunrise>
	keepers.AppKeepers
	// <sunrise>
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry
	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	// <sunrise>
	// the module manager
	mm *module.Manager

	// module configurator
	configurator module.Configurator
	// </sunrise>

	// simulation manager
	sm *module.SimulationManager

	// custom structure for skip-mev protection
	MevLane        *mevlane.MEVLane
	CheckTxHandler checktx.CheckTx
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)
}

// getGovProposalHandlers return the chain proposal handlers.
func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	// this line is used by starport scaffolding # stargate/app/govProposalHandlers

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		// this line is used by starport scaffolding # stargate/app/govProposalHandler
	)

	return govProposalHandlers
}

// AppConfig returns the default app config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		appConfig,
		// Loads the ao config from a YAML file.
		// appconfig.LoadYAML(AppConfigYAML),
		depinject.Supply(
			// supply custom module basics
			map[string]module.AppModuleBasic{
				genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
				// govtypes.ModuleName:     gov.NewAppModuleBasic(getGovProposalHandlers()),

				// overrides
				banktypes.ModuleName:   defaultoverrides.BankModuleBasic{},
				crisistypes.ModuleName: defaultoverrides.CrisisModuleBasic{},
				govtypes.ModuleName: defaultoverrides.GovModuleBasic{
					AppModuleBasic: gov.NewAppModuleBasic(getGovProposalHandlers()),
				},
				minttypes.ModuleName:           defaultoverrides.MintModuleBasic{},
				stakingtypes.ModuleName:        defaultoverrides.StakingModuleBasic{},
				tokenconvertertypes.ModuleName: defaultoverrides.TokenConverterModuleBasic{},
				feetypes.ModuleName:            defaultoverrides.FeeModuleBasic{},

				// this line is used by starport scaffolding # stargate/appConfig/moduleBasic
			},
		),
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
) (*App, error) {
	var (
		app        = &App{}
		appBuilder *runtime.AppBuilder

		// merge the AppConfig and other configuration in one config
		appConfig = depinject.Configs(
			AppConfig(),
			depinject.Supply(
				// Supply the application options
				appOpts,
				// Supply with IBC keeper getter for the IBC modules with App Wiring.
				// The IBC Keeper cannot be passed because it has not been initiated yet.
				// Passing the getter, the app IBC Keeper will always be accessible.
				// This needs to be removed after IBC supports App Wiring.
				app.GetIBCKeeper,
				app.GetCapabilityScopedKeeper,
				// Supply the logger
				logger,

				// ADVANCED CONFIGURATION
				//
				// AUTH
				//
				// For providing a custom function required in auth to generate custom account types
				// add it below. By default the auth module uses simulation.RandomGenesisAccounts.
				//
				// authtypes.RandomGenesisAccountsFn(simulation.RandomGenesisAccounts),
				//
				// For providing a custom a base account type add it below.
				// By default the auth module uses authtypes.ProtoBaseAccount().
				//
				// func() sdk.AccountI { return authtypes.ProtoBaseAccount() },
				//
				// For providing a different address codec, add it below.
				// By default the auth module uses a Bech32 address codec,
				// with the prefix defined in the auth module configuration.
				//
				// func() address.Codec { return <- custom address codec type -> }

				//
				// STAKING
				//
				// For provinding a different validator and consensus address codec, add it below.
				// By default the staking module uses the bech32 prefix provided in the auth config,
				// and appends "valoper" and "valcons" for validator and consensus addresses respectively.
				// When providing a custom address codec in auth, custom address codecs must be provided here as well.
				//
				// func() runtime.ValidatorAddressCodec { return <- custom validator address codec type -> }
				// func() runtime.ConsensusAddressCodec { return <- custom consensus address codec type -> }
			),
			depinject.Provide(
				//
				// MINT
				//

				// For providing a custom inflation function for x/mint add here your
				// custom function that implements the minttypes.InflationCalculationFn
				// interface.
				ProvideInflationCalculatorFn,
			),
		)
	)

	if err := depinject.Inject(appConfig,
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AppKeepers.AccountKeeper,
		&app.AppKeepers.BankKeeper,
		&app.AppKeepers.StakingKeeper,
		&app.AppKeepers.SlashingKeeper,
		&app.AppKeepers.MintKeeper,
		&app.AppKeepers.DistrKeeper,
		&app.AppKeepers.GovKeeper,
		&app.AppKeepers.CrisisKeeper,
		&app.AppKeepers.UpgradeKeeper,
		&app.AppKeepers.ParamsKeeper,
		&app.AppKeepers.AuthzKeeper,
		&app.AppKeepers.EvidenceKeeper,
		&app.AppKeepers.FeeGrantKeeper,
		&app.AppKeepers.GroupKeeper,
		&app.AppKeepers.ConsensusParamsKeeper,
		&app.AppKeepers.CircuitBreakerKeeper,

		// Third party module keepers
		&app.AppKeepers.AuctionKeeper,
		&app.AppKeepers.BlobKeeper,
		&app.AppKeepers.StreamKeeper,
		&app.AppKeepers.TokenconverterKeeper,
		&app.AppKeepers.LiquiditypoolKeeper,
		&app.AppKeepers.LiquidityincentiveKeeper,
		&app.AppKeepers.SwapKeeper,
		&app.AppKeepers.FeeKeeper,
		// this line is used by starport scaffolding # stargate/app/keeperDefinition
	); err != nil {
		panic(err)
	}

	// Below we could construct and set an application specific mempool and
	// ABCI 1.0 PrepareProposal and ProcessProposal handlers. These defaults are
	// already set in the SDK's BaseApp, this shows an example of how to override
	// them.
	//
	// Example:
	//
	// app.App = appBuilder.Build(...)
	// nonceMempool := mempool.NewSenderNonceMempool()
	// abciPropHandler := NewDefaultProposalHandler(nonceMempool, app.App.BaseApp)
	//
	// app.App.BaseApp.SetMempool(nonceMempool)
	// app.App.BaseApp.SetPrepareProposal(abciPropHandler.PrepareProposalHandler())
	// app.App.BaseApp.SetProcessProposal(abciPropHandler.ProcessProposalHandler())
	//
	// Alternatively, you can construct BaseApp options, append those to
	// baseAppOptions and pass them to the appBuilder.
	//
	// Example:
	//
	// prepareOpt = func(app *baseapp.BaseApp) {
	// 	abciPropHandler := baseapp.NewDefaultProposalHandler(nonceMempool, app)
	// 	app.SetPrepareProposal(abciPropHandler.PrepareProposalHandler())
	// }
	// baseAppOptions = append(baseAppOptions, prepareOpt)
	//
	// create and set vote extension handler
	// voteExtOp := func(bApp *baseapp.BaseApp) {
	// 	voteExtHandler := NewVoteExtensionHandler()
	// 	voteExtHandler.SetHandlers(bApp)
	// }

	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// Register legacy modules
	app.registerIBCModules()

	// <sunrise>
	app.AppKeepers.SwapKeeper.TransferKeeper = &app.AppKeepers.TransferKeeper
	// </sunrise>

	// register streaming services
	if err := app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
		return nil, err
	}

	/****  Module Options ****/

	// ---------------------------------------------------------------------------- //
	// ------------------------- Begin `Skip MEV` Code ---------------------------- //
	// ---------------------------------------------------------------------------- //
	// STEP 1-3: Create the Block SDK lanes.
	mevLane, freeLane, defaultLane := CreateLanes(app)

	// STEP 4: Construct a mempool based off the lanes. Note that the order of the lanes
	// matters. Blocks are constructed from the top lane to the bottom lane. The top lane
	// is the first lane in the array and the bottom lane is the last lane in the array.
	mempool, err := block.NewLanedMempool(
		app.Logger(),
		[]block.Lane{mevLane, freeLane, defaultLane},
	)
	if err != nil {
		panic(err)
	}

	// The application's mempool is now powered by the Block SDK!
	app.App.SetMempool(mempool)

	// STEP 5: Create a global ante handler that will be called on each transaction when
	// proposals are being built and verified. Note that this step must be done before
	// setting the ante handler on the lanes.
	anteHandler := ante.NewAnteHandler(
		app.AppKeepers.AccountKeeper,
		app.AppKeepers.BankKeeper,
		app.AppKeepers.FeeGrantKeeper,
		app.AppKeepers.BlobKeeper,
		app.AppKeepers.FeeKeeper,
		app.txConfig.SignModeHandler(),
		ante.DefaultSigVerificationGasConsumer,
		app.AppKeepers.IBCKeeper,
		app.AppKeepers.AuctionKeeper,
		mevLane,
		app.txConfig.TxEncoder(),
	)
	// Set the ante handler on the lanes.
	opt := []base.LaneOption{
		base.WithAnteHandler(anteHandler),
	}
	mevLane.WithOptions(
		opt...,
	)
	freeLane.WithOptions(
		opt...,
	)
	defaultLane.WithOptions(
		opt...,
	)

	app.MevLane = mevLane

	// Step 6: Create the proposal handler and set it on the app. Now the application
	// will build and verify proposals using the Block SDK!
	proposalHandler := abci.NewProposalHandler(
		app.Logger(),
		app.txConfig.TxDecoder(),
		app.txConfig.TxEncoder(),
		mempool,
	)
	app.App.SetPrepareProposal(proposalHandler.PrepareProposalHandler())
	app.App.SetProcessProposal(proposalHandler.ProcessProposalHandler())

	// Step 7: Set the custom CheckTx handler on BaseApp. This is only required if you
	// use the MEV lane.
	mevCheckTx := checktx.NewMEVCheckTxHandler(
		app.App,
		app.txConfig.TxDecoder(),
		mevLane,
		anteHandler,
		app.App.CheckTx,
	)
	checkTxHandler := checktx.NewMempoolParityCheckTx(
		app.Logger(), mempool,
		app.txConfig.TxDecoder(), mevCheckTx.CheckTx(),
	)

	app.SetCheckTx(checkTxHandler.CheckTx())

	// <sunrise>
	// Step 8: Set the custom Upgrade handler on BaseApp. This is added for on-chain upgrade.
	app.SetupUpgradeHandlers()
	// Step 8: Set the custom upgrade store loaders on BaseApp.
	app.SetupUpgradeStoreLoaders()
	// </sunrise>

	// ---------------------------------------------------------------------------- //
	// ------------------------- End `Skip MEV` Code ------------------------------ //
	// ---------------------------------------------------------------------------- //

	app.ModuleManager.RegisterInvariants(app.AppKeepers.CrisisKeeper)

	// add test gRPC service for testing gRPC queries in isolation
	testdata_pulsar.RegisterQueryServer(app.GRPCQueryRouter(), testdata_pulsar.QueryImpl{})

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AppKeepers.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, overrideModules)

	app.sm.RegisterStoreDecoders()

	// A custom InitChainer can be set if extra pre-init-genesis logic is required.
	// By default, when using app wiring enabled module, this is not required.
	// For instance, the upgrade module will set automatically the module version map in its init genesis thanks to app wiring.
	// However, when registering a module manually (i.e. that does not support app wiring), the module version map
	// must be set manually as follow. The upgrade module will de-duplicate the module version map.
	//
	// app.SetInitChainer(func(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	// 	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap())
	// 	return app.App.InitChainer(ctx, req)
	// })

	if err := app.Load(loadLatest); err != nil {
		return nil, err
	}

	return app, nil
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

// GetKey returns the KVStoreKey for the provided store key.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	kvStoreKey, ok := app.UnsafeFindStoreKey(storeKey).(*storetypes.KVStoreKey)
	if !ok {
		return nil
	}
	return kvStoreKey
}

// GetMemKey returns the MemoryStoreKey for the provided store key.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	key, ok := app.UnsafeFindStoreKey(storeKey).(*storetypes.MemoryStoreKey)
	if !ok {
		return nil
	}

	return key
}

// kvStoreKeys returns all the kv store keys registered inside App.
func (app *App) kvStoreKeys() map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storetypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	return keys
}

// GetSubspace returns a param subspace for a given module name.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.AppKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface.
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	// Register the base app API routes.
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)

	// Register the Block SDK mempool API routes.
	service.RegisterGRPCGatewayRoutes(apiSvr.ClientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API in app.go so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}

	// register app's OpenAPI routes.
	docs.RegisterOpenAPIService(Name, apiSvr.Router)
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	// Register the base app transaction service.
	app.App.RegisterTxService(clientCtx)

	// Register the Block SDK mempool transaction service.
	mempool, ok := app.App.Mempool().(block.Mempool)
	if !ok {
		panic("mempool is not a block.Mempool")
	}
	service.RegisterMempoolService(app.GRPCQueryRouter(), mempool)
}

// GetIBCKeeper returns the IBC keeper.
func (app *App) GetIBCKeeper() *ibckeeper.Keeper {
	return app.AppKeepers.IBCKeeper
}

// GetCapabilityScopedKeeper returns the capability scoped keeper.
func (app *App) GetCapabilityScopedKeeper(moduleName string) capabilitykeeper.ScopedKeeper {
	return app.AppKeepers.CapabilityKeeper.ScopeToModule(moduleName)
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

// GetTxConfig implements the TestingApp interface.
func (app *App) GetTxConfig() client.TxConfig {
	return app.txConfig
}

func (app *App) SetupUpgradeStoreLoaders() {
	upgradeInfo, err := app.AppKeepers.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.AppKeepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		if upgradeInfo.Name == upgrade.UpgradeName {
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &upgrade.StoreUpgrades))
		}
	}

}

func (app *App) SetupUpgradeHandlers() {
	for _, upgrade := range Upgrades {
		app.AppKeepers.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				app.mm,
				app.configurator,
				app.BaseApp,
				&app.AppKeepers,
			),
		)
	}
}
