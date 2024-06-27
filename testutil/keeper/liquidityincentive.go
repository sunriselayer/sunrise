package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	liquidityincentivetestutil "github.com/sunriselayer/sunrise/x/liquidityincentive/testutil"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

type LiquidityIncentiveMocks struct {
	AcctKeeper          *liquidityincentivetestutil.MockAccountKeeper
	BankKeeper          *liquidityincentivetestutil.MockBankKeeper
	StakingKeeper       *liquidityincentivetestutil.MockStakingKeeper
	LiquiditypoolKeeper *liquidityincentivetestutil.MockLiquidityPoolKeeper
}

func LiquidityincentiveKeeper(t testing.TB) (keeper.Keeper, LiquidityIncentiveMocks, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	// gomock initializations
	ctrl := gomock.NewController(t)
	m := LiquidityIncentiveMocks{
		AcctKeeper:          liquidityincentivetestutil.NewMockAccountKeeper(ctrl),
		BankKeeper:          liquidityincentivetestutil.NewMockBankKeeper(ctrl),
		StakingKeeper:       liquidityincentivetestutil.NewMockStakingKeeper(ctrl),
		LiquiditypoolKeeper: liquidityincentivetestutil.NewMockLiquidityPoolKeeper(ctrl),
	}

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		m.AcctKeeper,
		m.BankKeeper,
		m.StakingKeeper,
		m.LiquiditypoolKeeper,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, m, ctx
}
