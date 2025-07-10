package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/lockup/keeper"
	module "github.com/sunriselayer/sunrise/x/lockup/module"
	lockuptestutil "github.com/sunriselayer/sunrise/x/lockup/testutil"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	mocks        LockupMocks
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx
	authority := authtypes.NewModuleAddress(types.GovModuleName)

	mocks := getMocks(t)

	k := keeper.NewKeeper(
		encCfg.Codec,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		addressCodec,
		mocks.AccountKeeper,
		mocks.BankKeeper,
		mocks.StakingKeeper,
		mocks.TokenConverterKeeper,
		mocks.ShareclassKeeper,
	)

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
		mocks:        mocks,
	}
}

type LockupMocks struct {
	AccountKeeper        *lockuptestutil.MockAccountKeeper
	BankKeeper           *lockuptestutil.MockBankKeeper
	StakingKeeper        *lockuptestutil.MockStakingKeeper
	TokenConverterKeeper *lockuptestutil.MockTokenConverterKeeper
	ShareclassKeeper     *lockuptestutil.MockShareclassKeeper
}

func getMocks(t *testing.T) LockupMocks {
	ctrl := gomock.NewController(t)

	return LockupMocks{
		AccountKeeper:        lockuptestutil.NewMockAccountKeeper(ctrl),
		BankKeeper:           lockuptestutil.NewMockBankKeeper(ctrl),
		StakingKeeper:        lockuptestutil.NewMockStakingKeeper(ctrl),
		TokenConverterKeeper: lockuptestutil.NewMockTokenConverterKeeper(ctrl),
		ShareclassKeeper:     lockuptestutil.NewMockShareclassKeeper(ctrl),
	}
}
