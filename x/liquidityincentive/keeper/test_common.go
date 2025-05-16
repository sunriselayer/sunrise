package keeper

import (
	"context"
	"testing"

	"cosmossdk.io/core/appmodule"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

type MockBankKeeper struct {
	types.BankKeeper
	balances map[string]sdk.Coins
}

func (k *MockBankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	if k.balances == nil {
		k.balances = make(map[string]sdk.Coins)
	}
	if _, exists := k.balances[moduleName]; !exists {
		k.balances[moduleName] = sdk.NewCoins()
	}
	k.balances[moduleName] = k.balances[moduleName].Add(amt...)
	return nil
}

func (k *MockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	if k.balances == nil {
		k.balances = make(map[string]sdk.Coins)
	}
	if _, exists := k.balances[senderModule]; !exists {
		k.balances[senderModule] = sdk.NewCoins()
	}
	if _, exists := k.balances[recipientAddr.String()]; !exists {
		k.balances[recipientAddr.String()] = sdk.NewCoins()
	}
	if !k.balances[senderModule].IsAllGTE(amt) {
		return errorsmod.Wrapf(errors.ErrInsufficientFunds, "module %s has insufficient funds: %s < %s", senderModule, k.balances[senderModule], amt)
	}
	k.balances[senderModule] = k.balances[senderModule].Sub(amt...)
	k.balances[recipientAddr.String()] = k.balances[recipientAddr.String()].Add(amt...)
	return nil
}

func (k *MockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	if k.balances == nil {
		k.balances = make(map[string]sdk.Coins)
	}
	if _, exists := k.balances[senderAddr.String()]; !exists {
		k.balances[senderAddr.String()] = sdk.NewCoins()
	}
	if _, exists := k.balances[recipientModule]; !exists {
		k.balances[recipientModule] = sdk.NewCoins()
	}
	if !k.balances[senderAddr.String()].IsAllGTE(amt) {
		return errorsmod.Wrapf(errors.ErrInsufficientFunds, "account %s has insufficient funds: %s < %s", senderAddr.String(), k.balances[senderAddr.String()], amt)
	}
	k.balances[senderAddr.String()] = k.balances[senderAddr.String()].Sub(amt...)
	k.balances[recipientModule] = k.balances[recipientModule].Add(amt...)
	return nil
}

func (k *MockBankKeeper) GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	if k.balances == nil {
		k.balances = make(map[string]sdk.Coins)
	}
	if _, exists := k.balances[addr.String()]; !exists {
		k.balances[addr.String()] = sdk.NewCoins()
	}
	coins := k.balances[addr.String()]
	return sdk.NewCoin(denom, coins.AmountOf(denom))
}

func (k *MockBankKeeper) IsSendEnabledCoins(ctx context.Context, coins ...sdk.Coin) error {
	return nil
}

type MockAccountKeeper struct {
	types.AccountKeeper
}

func (k *MockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return sdk.AccAddress([]byte(moduleName))
}

type MockStakingKeeper struct {
	types.StakingKeeper
}

type MockFeeKeeper struct {
	types.FeeKeeper
}

type MockTokenConverterKeeper struct {
	types.TokenConverterKeeper
}

type MockLiquidityPoolKeeper struct {
	types.LiquidityPoolKeeper
}

func setupKeeperWithParams(t *testing.T) (sdk.Context, Keeper) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	authority := authtypes.NewModuleAddress(types.ModuleName)

	env := appmodule.Environment{
		Logger:         log.NewNopLogger(),
		KVStoreService: runtime.NewKVStoreService(storeKey),
	}

	bankKeeper := &MockBankKeeper{}
	accountKeeper := &MockAccountKeeper{}
	stakingKeeper := &MockStakingKeeper{}
	feeKeeper := &MockFeeKeeper{}
	tokenConverterKeeper := &MockTokenConverterKeeper{}
	liquidityPoolKeeper := &MockLiquidityPoolKeeper{}

	// Create address codec for testing
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	k := NewKeeper(
		env,
		cdc,
		addressCodec,
		authority,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		feeKeeper,
		tokenConverterKeeper,
		liquidityPoolKeeper,
	)

	ctx := sdk.NewContext(stateStore, true, log.NewNopLogger())

	// Initialize params
	params := types.DefaultParams()
	err := k.Params.Set(ctx, params)
	require.NoError(t, err)

	return ctx, k
}
