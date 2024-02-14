package util

import (
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	// tmversion "github.com/cometbft/cometbft/proto/tendermint/version"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	"github.com/sunrise-zone/sunrise-app/app"
	"github.com/sunrise-zone/sunrise-app/testutil"

	"github.com/sunrise-zone/sunrise-app/test/util/testfactory"
	"github.com/sunrise-zone/sunrise-app/test/util/testnode"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const ChainID = testfactory.ChainID

// Get flags every time the simulator is run
func init() {
	simcli.GetSimulatorFlags()
}

type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(_ string) interface{} {
	return nil
}

// SetupTestAppWithGenesisValSet initializes a new app with a validator set and
// genesis accounts that also act as delegators. For simplicity, each validator
// is bonded with a delegation of one consensus engine unit in the default token
// of the app from first genesis account. A no-op logger is set in app.
func SetupTestAppWithGenesisValSet(cparams *tmproto.ConsensusParams, genAccounts ...string) (*app.App, keyring.Keyring) {
	testutil.InitSDKConfig()
	// var cache storetypes.MultiStorePersistentCache
	// EmptyAppOptions is a stub implementing AppOptions
	emptyOpts := EmptyAppOptions{}
	// var anteOpt = func(bapp *baseapp.BaseApp) { bapp.SetAnteHandler(nil) }
	db := dbm.NewMemDB()

	testApp, err := app.New(
		log.NewNopLogger(), db, nil, true,
		emptyOpts,
		func(bApp *baseapp.BaseApp) {
			bApp.FinalizeBlock(&abci.RequestFinalizeBlock{})
		},
	)

	if err != nil {
		panic(err)
	}

	genesisState, valSet, kr := GenesisStateWithSingleValidator(testApp, genAccounts...)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	genesisTime := time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC).UTC()

	// set chain id
	fn := baseapp.SetChainID(ChainID)
	fn(testApp.BaseApp)

	// init chain will set the validator set and initialize the genesis accounts
	_, err = testApp.InitChain(
		&abci.RequestInitChain{
			Time:            genesisTime,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: cparams,
			AppStateBytes:   stateBytes,
			ChainId:         ChainID,
		},
	)
	if err != nil {
		panic(err)
	}

	// commit genesis changes
	testApp.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:             testApp.LastBlockHeight() + 1,
		NextValidatorsHash: valSet.Hash(),
	})
	testApp.Commit()
	_ = valSet
	// testApp.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
	// 	ChainID:            ChainID,
	// 	Height:             testApp.LastBlockHeight() + 1,
	// 	AppHash:            testApp.LastCommitID().Hash,
	// 	ValidatorsHash:     valSet.Hash(),
	// 	NextValidatorsHash: valSet.Hash(),
	// 	Version: tmversion.Consensus{
	// 		App: cparams.Version.App,
	// 	},
	// }})

	return testApp, kr
}

// AddAccount mimics the cli addAccount command, providing an
// account with an allocation of to "token" and "tia" tokens in the genesis
// state
func AddAccount(addr sdk.AccAddress, appState app.GenesisState, cdc codec.Codec) (map[string]json.RawMessage, error) {
	// create concrete account type based on input parameters
	var genAccount authtypes.GenesisAccount

	coins := sdk.Coins{
		sdk.NewCoin("token", sdkmath.NewInt(1000000)),
		sdk.NewCoin(app.BondDenom, sdkmath.NewInt(1000000)),
	}

	balances := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
	baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

	genAccount = baseAccount

	if err := genAccount.Validate(); err != nil {
		return appState, fmt.Errorf("failed to validate new genesis account: %w", err)
	}

	authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return appState, fmt.Errorf("failed to get accounts from any: %w", err)
	}

	if accs.Contains(addr) {
		return appState, fmt.Errorf("cannot add account at existing address %s", addr)
	}

	// Add the new account to the set of genesis accounts and sanitize the
	// accounts afterwards.
	accs = append(accs, genAccount)
	accs = authtypes.SanitizeGenesisAccounts(accs)

	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		return appState, fmt.Errorf("failed to convert accounts into any's: %w", err)
	}
	authGenState.Accounts = genAccs

	authGenStateBz, err := cdc.MarshalJSON(&authGenState)
	if err != nil {
		return appState, fmt.Errorf("failed to marshal auth genesis state: %w", err)
	}

	appState[authtypes.ModuleName] = authGenStateBz

	bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
	bankGenState.Balances = append(bankGenState.Balances, balances)
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return appState, fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}

	appState[banktypes.ModuleName] = bankGenStateBz
	return appState, nil
}

// GenesisStateWithSingleValidator initializes GenesisState with a single
// validator and genesis accounts that also act as delegators.
func GenesisStateWithSingleValidator(testApp *app.App, genAccounts ...string) (app.GenesisState, *tmtypes.ValidatorSet, keyring.Keyring) {
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	if err != nil {
		panic(err)
	}

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	accs := make([]authtypes.GenesisAccount, 0, len(genAccounts)+1)
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	accs = append(accs, acc)
	balances := make([]banktypes.Balance, 0, len(genAccounts)+1)
	balances = append(balances, banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(app.BondDenom, sdkmath.NewInt(100000000000000))),
	})

	kr, fundedBankAccs, fundedAuthAccs := testnode.FundKeyringAccounts(genAccounts...)
	accs = append(accs, fundedAuthAccs...)
	balances = append(balances, fundedBankAccs...)

	genesisState := NewDefaultGenesisState(testApp.AppCodec())
	genesisState = genesisStateWithValSet(testApp, genesisState, valSet, accs, balances...)

	return genesisState, valSet, kr
}

func genesisStateWithValSet(
	a *app.App,
	genesisState app.GenesisState,
	valSet *tmtypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) app.GenesisState {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = a.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		if err != nil {
			panic(err)
		}
		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			panic(err)
		}
		validator := stakingtypes.Validator{
			OperatorAddress: sdk.ValAddress(val.Address).String(),
			ConsensusPubkey: pkAny,
			Jailed:          false,
			Status:          stakingtypes.Bonded,
			Tokens:          bondAmt,
			DelegatorShares: sdkmath.LegacyOneDec(),
			Description: stakingtypes.Description{
				Moniker:         "test-moniker",
				Identity:        "test-identity",
				Website:         "https://www.google.com/",
				SecurityContact: "sunrise17p9rzwnnfxcjp32un9ug7yhhzgtkhvl9jfksztgw5uh69wac2pgs06edvm",
				Details:         "test-details",
			},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()),
			MinSelfDelegation: sdkmath.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(val.Address).String(), sdkmath.LegacyOneDec()))

	}
	// set validators and delegations
	params := stakingtypes.DefaultParams()
	params.BondDenom = app.BondDenom
	stakingGenesis := stakingtypes.NewGenesisState(params, validators, delegations)
	genesisState[stakingtypes.ModuleName] = a.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(app.BondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(app.BondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = a.AppCodec().MustMarshalJSON(bankGenesis)

	return genesisState
}

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) app.GenesisState {
	return app.ModuleBasics().DefaultGenesis(cdc)
}
