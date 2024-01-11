package util

import (
	"bytes"
	"testing"
	"time"

	blobstream "github.com/sunrise-zone/sunrise-app/x/blobstream/module"

	"github.com/sunrise-zone/sunrise-app/app"
	"github.com/sunrise-zone/sunrise-app/x/blobstream/keeper"
	bstypes "github.com/sunrise-zone/sunrise-app/x/blobstream/types"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ccodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ccrypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	// ModuleBasics is a mock module basic manager for testing
	ModuleBasics = app.ModuleBasics
	// TestingStakeParams is a set of staking params for testing
	TestingStakeParams = stakingtypes.Params{
		UnbondingTime:     100,
		MaxValidators:     10,
		MaxEntries:        10,
		HistoricalEntries: 10000,
		BondDenom:         "stake",
		MinCommissionRate: sdkmath.LegacyNewDecWithPrec(0, 0),
	}

	// ConsPrivKeys generate ed25519 ConsPrivKeys to be used for validator operator keys
	ConsPrivKeys = []ccrypto.PrivKey{
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
	}

	// ConsPubKeys holds the consensus public keys to be used for validator operator keys
	ConsPubKeys = []ccrypto.PubKey{
		ConsPrivKeys[0].PubKey(),
		ConsPrivKeys[1].PubKey(),
		ConsPrivKeys[2].PubKey(),
		ConsPrivKeys[3].PubKey(),
		ConsPrivKeys[4].PubKey(),
	}

	// AccPrivKeys generate secp256k1 pubkeys to be used for account pub keys
	AccPrivKeys = []ccrypto.PrivKey{
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
	}

	// AccPubKeys holds the pub keys for the account keys
	AccPubKeys = []ccrypto.PubKey{
		AccPrivKeys[0].PubKey(),
		AccPrivKeys[1].PubKey(),
		AccPrivKeys[2].PubKey(),
		AccPrivKeys[3].PubKey(),
		AccPrivKeys[4].PubKey(),
	}

	// AccAddrs holds the sdk.AccAddresses
	AccAddrs = []sdk.AccAddress{
		sdk.AccAddress(AccPubKeys[0].Address()),
		sdk.AccAddress(AccPubKeys[1].Address()),
		sdk.AccAddress(AccPubKeys[2].Address()),
		sdk.AccAddress(AccPubKeys[3].Address()),
		sdk.AccAddress(AccPubKeys[4].Address()),
	}

	// ValAddrs holds the sdk.ValAddresses
	ValAddrs = []sdk.ValAddress{
		sdk.ValAddress(AccPubKeys[0].Address()),
		sdk.ValAddress(AccPubKeys[1].Address()),
		sdk.ValAddress(AccPubKeys[2].Address()),
		sdk.ValAddress(AccPubKeys[3].Address()),
		sdk.ValAddress(AccPubKeys[4].Address()),
	}

	// EVMAddrs holds etheruem addresses
	EVMAddrs = initEVMAddrs(100)

	// InitTokens holds the number of tokens to initialize an account with
	InitTokens = sdk.TokensFromConsensusPower(110, sdk.DefaultPowerReduction)

	// InitCoins holds the number of coins to initialize an account with
	InitCoins = sdk.NewCoins(sdk.NewCoin(TestingStakeParams.BondDenom, InitTokens))

	// StakingAmount holds the staking power to start a validator with
	StakingAmount = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
)

func initEVMAddrs(count int) []gethcommon.Address {
	addresses := make([]gethcommon.Address, count)
	for i := 0; i < count; i++ {
		evmAddr := gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(i + 1)}, gethcommon.AddressLength))
		addresses[i] = evmAddr
	}
	return addresses
}

// TestInput stores the various keepers required to test Blobstream
type TestInput struct {
	BlobstreamKeeper keeper.Keeper
	AccountKeeper    authkeeper.AccountKeeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	DistKeeper       distrkeeper.Keeper
	BankKeeper       bankkeeper.Keeper
	Context          sdk.Context
	Marshaler        codec.Codec
	LegacyAmino      *codec.LegacyAmino
}

// CreateTestEnvWithoutBlobstreamKeysInit creates the keeper testing environment for Blobstream
func CreateTestEnvWithoutBlobstreamKeysInit(t *testing.T) TestInput {
	t.Helper()

	db := dbm.NewMemDB()
	emptyOpts := EmptyAppOptions{}
	// encCfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)

	testApp, _ := app.New(
		log.NewNopLogger(), db, nil, true,
		emptyOpts,
	)

	return TestInput{
		BlobstreamKeeper: testApp.StreamKeeper,
		AccountKeeper:    testApp.AccountKeeper,
		BankKeeper:       testApp.BankKeeper,
		StakingKeeper:    *testApp.StakingKeeper,
		SlashingKeeper:   testApp.SlashingKeeper,
		DistKeeper:       testApp.DistrKeeper,
		Context:          testApp.NewContext(false),
		Marshaler:        testApp.AppCodec(),
		LegacyAmino:      testApp.LegacyAmino(),
	}
}

// CreateTestEnv creates the keeper testing environment for Blobstream
func CreateTestEnv(t *testing.T) TestInput {
	input := CreateTestEnvWithoutBlobstreamKeysInit(t)
	input.BlobstreamKeeper.SetLatestAttestationNonce(input.Context, blobstream.InitialLatestAttestationNonce)
	input.BlobstreamKeeper.SetEarliestAvailableAttestationNonce(input.Context, blobstream.InitialEarliestAvailableAttestationNonce)
	return input
}

// MakeTestCodec creates a legacy amino codec for testing
func MakeTestCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	auth.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	bank.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	staking.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	distribution.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	sdk.RegisterLegacyAminoCodec(cdc)
	ccodec.RegisterCrypto(cdc)
	params.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	return cdc
}

// getSubspace returns a param subspace for a given module name.
func getSubspace(k paramskeeper.Keeper, moduleName string) paramstypes.Subspace {
	subspace, _ := k.GetSubspace(moduleName)
	return subspace
}

// MakeTestMarshaler creates a proto codec for use in testing
func MakeTestMarshaler() codec.Codec {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	bstypes.RegisterInterfaces(interfaceRegistry)
	return codec.NewProtoCodec(interfaceRegistry)
}

// SetupFiveValChain does all the initialization for a 5 Validator chain using the keys here
func SetupFiveValChain(t *testing.T) (TestInput, sdk.Context) {
	t.Helper()
	input := CreateTestEnv(t)

	// Set the params for our modules
	input.StakingKeeper.SetParams(input.Context, TestingStakeParams)

	// Initialize each of the validators
	for i := range []int{0, 1, 2, 3, 4} {
		CreateValidator(t, input, AccAddrs[i], AccPubKeys[i], uint64(i), ValAddrs[i], ConsPubKeys[i], StakingAmount)
		RegisterEVMAddress(t, input, ValAddrs[i], EVMAddrs[i])
	}

	// Run the staking endblocker to ensure valset is correct in state
	input.StakingKeeper.EndBlocker(input.Context)

	// Return the test input
	return input, input.Context
}

func CreateValidator(
	t *testing.T,
	input TestInput,
	accAddr sdk.AccAddress,
	accPubKey ccrypto.PubKey,
	accountNumber uint64,
	valAddr sdk.ValAddress,
	consPubKey ccrypto.PubKey,
	stakingAmount sdkmath.Int,
) {
	// Initialize the account for the key
	acc := input.AccountKeeper.NewAccount(
		input.Context,
		authtypes.NewBaseAccount(accAddr, accPubKey, accountNumber, 0),
	)

	// Set the balance for the account
	require.NoError(t, input.BankKeeper.MintCoins(input.Context, bstypes.ModuleName, InitCoins))
	err := input.BankKeeper.SendCoinsFromModuleToAccount(input.Context, bstypes.ModuleName, acc.GetAddress(), InitCoins)
	require.NoError(t, err)

	// Set the account in state
	input.AccountKeeper.SetAccount(input.Context, acc)

	// Create a validator for that account using some tokens in the account
	// and the staking handler
	msgServer := stakingkeeper.NewMsgServerImpl(&input.StakingKeeper)
	_, err = msgServer.CreateValidator(input.Context, NewTestMsgCreateValidator(valAddr, consPubKey, stakingAmount))
	require.NoError(t, err)
}

func RegisterEVMAddress(
	t *testing.T,
	input TestInput,
	valAddr sdk.ValAddress,
	evmAddr gethcommon.Address,
) {
	bsMsgServer := keeper.NewMsgServerImpl(input.BlobstreamKeeper)
	registerMsg := bstypes.NewMsgRegisterEvmAddress(valAddr, evmAddr)
	_, err := bsMsgServer.RegisterEvmAddress(input.Context, registerMsg)
	require.NoError(t, err)
}

func NewTestMsgCreateValidator(
	address sdk.ValAddress,
	pubKey ccrypto.PubKey,
	amt sdkmath.Int,
) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec())
	out, err := stakingtypes.NewMsgCreateValidator(
		address.String(), pubKey, sdk.NewCoin("stake", amt),
		stakingtypes.Description{
			Moniker:         "",
			Identity:        "",
			Website:         "",
			SecurityContact: "",
			Details:         "",
		}, commission, sdkmath.OneInt(),
	)
	if err != nil {
		panic(err)
	}
	return out
}

// SetupTestChain sets up a test environment with the provided validator voting weights
func SetupTestChain(t *testing.T, weights []uint64) (TestInput, sdk.Context) {
	t.Helper()
	input := CreateTestEnv(t)

	// Set the params for our modules
	TestingStakeParams.MaxValidators = 100
	input.StakingKeeper.SetParams(input.Context, TestingStakeParams)

	// Initialize each of the validators
	stakingMsgServer := stakingkeeper.NewMsgServerImpl(&input.StakingKeeper)
	bsMsgServer := keeper.NewMsgServerImpl(input.BlobstreamKeeper)
	for i, weight := range weights {
		consPrivKey := ed25519.GenPrivKey()
		consPubKey := consPrivKey.PubKey()
		valPrivKey := secp256k1.GenPrivKey()
		valPubKey := valPrivKey.PubKey()
		valAddr := sdk.ValAddress(valPubKey.Address())
		accAddr := sdk.AccAddress(valPubKey.Address())

		// Initialize the account for the key
		acc := input.AccountKeeper.NewAccount(
			input.Context,
			authtypes.NewBaseAccount(accAddr, valPubKey, uint64(i), 0),
		)

		// Set the balance for the account
		weightCoins := sdk.NewCoins(sdk.NewInt64Coin(TestingStakeParams.BondDenom, int64(weight)))
		require.NoError(t, input.BankKeeper.MintCoins(input.Context, bstypes.ModuleName, weightCoins))
		require.NoError(t, input.BankKeeper.SendCoinsFromModuleToAccount(input.Context, bstypes.ModuleName, accAddr, weightCoins))

		// Set the account in state
		input.AccountKeeper.SetAccount(input.Context, acc)

		// Create a validator for that account using some of the tokens in the account
		// and the staking handler
		_, err := stakingMsgServer.CreateValidator(
			input.Context,
			NewTestMsgCreateValidator(valAddr, consPubKey, sdkmath.NewIntFromUint64(weight)),
		)
		require.NoError(t, err)

		registerMsg := bstypes.NewMsgRegisterEvmAddress(valAddr, EVMAddrs[i])
		_, err = bsMsgServer.RegisterEvmAddress(input.Context, registerMsg)
		require.NoError(t, err)

		// Run the staking endblocker to ensure valset is correct in state
		input.StakingKeeper.EndBlocker(input.Context)
	}

	// some inputs can cause the validator creation not to work, this checks that
	// everything was successful
	validators, _ := input.StakingKeeper.GetBondedValidatorsByPower(input.Context)
	require.Equal(t, len(weights), len(validators))

	// Return the test input
	return input, input.Context
}

func NewTestMsgUnDelegateValidator(address sdk.ValAddress, amt sdkmath.Int) *stakingtypes.MsgUndelegate {
	msg := stakingtypes.NewMsgUndelegate(sdk.AccAddress(address).String(), address.String(), sdk.NewCoin("stake", amt))
	return msg
}

// ExecuteBlobstreamHeights executes the end exclusive range of heights specified by beginHeight and endHeight
// along with the Blobstream abci.EndBlocker on each one of them.
// Returns the updated context with block height advanced to endHeight.
func ExecuteBlobstreamHeights(ctx sdk.Context, bsKeeper keeper.Keeper, beginHeight int64, endHeight int64) sdk.Context {
	for i := beginHeight; i < endHeight; i++ {
		ctx = ctx.WithBlockHeight(i)
		blobstream.EndBlocker(ctx, bsKeeper)
	}
	return ctx
}

// ExecuteBlobstreamHeightsWithTime executes the end exclusive range of heights specified by beginHeight and endHeight
// along with the Blobstream abci.EndBlocker on each one of them.
// Uses the interval to calculate the block header time.
func ExecuteBlobstreamHeightsWithTime(ctx sdk.Context, bsKeeper keeper.Keeper, beginHeight int64, endHeight int64, blockInterval time.Duration) sdk.Context {
	blockTime := ctx.BlockTime()
	for i := beginHeight; i < endHeight; i++ {
		ctx = ctx.WithBlockHeight(i).WithBlockTime(blockTime)
		blobstream.EndBlocker(ctx, bsKeeper)
		blockTime = blockTime.Add(blockInterval)
	}
	return ctx
}
