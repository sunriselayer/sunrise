package app

import (
	"fmt"
	"io"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/registry"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/accounts"
	authzkeeper "cosmossdk.io/x/authz/keeper"
	bankkeeper "cosmossdk.io/x/bank/keeper"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	consensuskeeper "cosmossdk.io/x/consensus/keeper"
	distrkeeper "cosmossdk.io/x/distribution/keeper"
	epochskeeper "cosmossdk.io/x/epochs/keeper"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	govkeeper "cosmossdk.io/x/gov/keeper"
	groupkeeper "cosmossdk.io/x/group/keeper"
	mintkeeper "cosmossdk.io/x/mint/keeper"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	_ "cosmossdk.io/x/protocolpool"
	poolkeeper "cosmossdk.io/x/protocolpool/keeper"
	slashingkeeper "cosmossdk.io/x/slashing/keeper"
	stakingkeeper "cosmossdk.io/x/staking/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante/unorderedtx"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	abci "github.com/cometbft/cometbft/abci/types"
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

	paramskeeper "cosmossdk.io/x/params/keeper"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v9/modules/apps/27-interchain-accounts/controller/keeper"
	icahostkeeper "github.com/cosmos/ibc-go/v9/modules/apps/27-interchain-accounts/host/keeper"
	ibcfeekeeper "github.com/cosmos/ibc-go/v9/modules/apps/29-fee/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v9/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v9/modules/core/keeper"

	"github.com/sunriselayer/sunrise/docs"
	damodulekeeper "github.com/sunriselayer/sunrise/x/da/keeper"
	feemodulekeeper "github.com/sunriselayer/sunrise/x/fee/keeper"
	liquidityincentivemodulekeeper "github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	liquiditypoolmodulekeeper "github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	swapmodulekeeper "github.com/sunriselayer/sunrise/x/swap/keeper"
	tokenconvertermodulekeeper "github.com/sunriselayer/sunrise/x/tokenconverter/keeper"

	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/sunriselayer/sunrise/app/mint"

	basedepinject "cosmossdk.io/x/accounts/defaults/base/depinject"
	lockupdepinject "cosmossdk.io/x/accounts/defaults/lockup/depinject"
	multisigdepinject "cosmossdk.io/x/accounts/defaults/multisig/depinject"

	selfdelegatablelockupdepinject "github.com/sunriselayer/sunrise/x/accounts/self_delegatable_lockup/depinject"
	selfdelegationproxydepinject "github.com/sunriselayer/sunrise/x/accounts/self_delegation_proxy/depinject"
	selfdelegationmodulekeeper "github.com/sunriselayer/sunrise/x/selfdelegation/keeper"
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
	legacyAmino       registry.AminoRegistrar
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// keepers
	AccountsKeeper        accounts.Keeper
	AuthKeeper            authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            *mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             *govkeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	NFTKeeper             nftkeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper
	CircuitBreakerKeeper  circuitkeeper.Keeper
	PoolKeeper            poolkeeper.Keeper
	EpochsKeeper          *epochskeeper.Keeper

	// IBC
	ParamsKeeper        paramskeeper.Keeper
	IBCKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper        ibcfeekeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper

	DaKeeper                 damodulekeeper.Keeper
	TokenconverterKeeper     tokenconvertermodulekeeper.Keeper
	LiquiditypoolKeeper      liquiditypoolmodulekeeper.Keeper
	LiquidityincentiveKeeper liquidityincentivemodulekeeper.Keeper
	SwapKeeper               swapmodulekeeper.Keeper
	FeeKeeper                feemodulekeeper.Keeper
	SelfdelegationKeeper     selfdelegationmodulekeeper.Keeper
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
		depinject.Provide(mint.ProvideMintFn),
	)
}

// New returns a reference to an initialized App.
func New(
	logger log.Logger,
	db corestore.KVStoreWithBatch,
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
			depinject.Provide(
				basedepinject.ProvideAccount,
				lockupdepinject.ProvideAllLockupAccounts,
				multisigdepinject.ProvideAccount,

				basedepinject.ProvideSecp256K1PubKey,

				selfdelegatablelockupdepinject.ProvideAllLockupAccounts,
				selfdelegationproxydepinject.ProvideAccount,
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
		&app.AccountsKeeper,
		&app.BankKeeper,
		&app.StakingKeeper,
		&app.SlashingKeeper,
		&app.MintKeeper,
		&app.DistrKeeper,
		&app.GovKeeper,
		&app.UpgradeKeeper,
		&app.AuthzKeeper,
		&app.EvidenceKeeper,
		&app.FeeGrantKeeper,
		&app.GroupKeeper,
		&app.NFTKeeper,
		&app.ConsensusParamsKeeper,
		&app.CircuitBreakerKeeper,
		&app.PoolKeeper,
		&app.EpochsKeeper,
		&app.DaKeeper,
		&app.TokenconverterKeeper,
		&app.LiquiditypoolKeeper,
		&app.LiquidityincentiveKeeper,
		&app.SwapKeeper,
		&app.FeeKeeper,
		&app.SelfdelegationKeeper,
	); err != nil {
		panic(err)
	}

	// add to default baseapp options
	// enable optimistic execution
	baseAppOptions = append(baseAppOptions, baseapp.SetOptimisticExecution())

	// build app
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// Register legacy modules
	app.registerIBCModules()

	// <sunrise>
	app.SwapKeeper.TransferKeeper = &app.TransferKeeper
	// </sunrise>

	// register streaming services
	if err := app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
		panic(err)
	}

	/****  Module Options ****/

	// <sunrise>
	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			ante.HandlerOptions{
				AccountKeeper:            app.AuthKeeper,
				BankKeeper:               app.BankKeeper,
				ConsensusKeeper:          app.ConsensusParamsKeeper,
				SignModeHandler:          app.txConfig.SignModeHandler(),
				FeegrantKeeper:           app.FeeGrantKeeper,
				SigGasConsumer:           ante.DefaultSigVerificationGasConsumer,
				UnorderedTxManager:       app.UnorderedTxManager,
				Environment:              app.AuthKeeper.Environment,
				AccountAbstractionKeeper: app.AccountsKeeper,
			},
			&app.CircuitBreakerKeeper,
			app.FeeKeeper,
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
		baseapp.NewDefaultProposalHandler(app.Mempool(), app),
	)

	app.BaseApp.SetProcessProposal(daProposalHandler.ProcessProposal())
	// </sunrise>

	// create the simulation manager and define the order of the modules for deterministic simulations
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AuthKeeper, &app.AccountsKeeper, authsims.RandomGenesisAccounts, nil),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, overrideModules)

	app.sm.RegisterStoreDecoders()

	// A custom InitChainer sets if extra pre-init-genesis logic is required.
	// This is necessary for manually registered modules that do not support app wiring.
	// Manually set the module version map as shown below.
	// The upgrade module will automatically handle de-duplication of the module version map.
	app.SetInitChainer(func(ctx sdk.Context, req *abci.InitChainRequest) (*abci.InitChainResponse, error) {
		if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap()); err != nil {
			return nil, err
		}
		return app.App.InitChainer(ctx, req)
	})

	// register custom snapshot extensions (if any)
	if manager := app.SnapshotManager(); manager != nil {
		if err := manager.RegisterExtensions(
			unorderedtx.NewSnapshotter(app.UnorderedTxManager),
		); err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %w", err))
		}
	}

	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	return app
}

// LegacyAmino returns App's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	switch cdc := app.legacyAmino.(type) {
	case *codec.LegacyAmino:
		return cdc
	default:
		panic("unexpected codec type")
	}
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
