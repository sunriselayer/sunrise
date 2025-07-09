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
	"github.com/sunriselayer/sunrise/x/lockup/testutil/mocks"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

type LockupMocks struct {
	AccountKeeper        *mocks.MockAccountKeeper
	BankKeeper           *mocks.MockBankKeeper
	StakingKeeper        *mocks.MockStakingKeeper
	TokenConverterKeeper *mocks.MockTokenConverterKeeper
	ShareclassKeeper     *mocks.MockShareclassKeeper
}
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

	ctrl := gomock.NewController(t)
	mockKeepers := LockupMocks{
		AccountKeeper:        mocks.NewMockAccountKeeper(ctrl),
		BankKeeper:           mocks.NewMockBankKeeper(ctrl),
		StakingKeeper:        mocks.NewMockStakingKeeper(ctrl),
		TokenConverterKeeper: mocks.NewMockTokenConverterKeeper(ctrl),
		ShareclassKeeper:     mocks.NewMockShareclassKeeper(ctrl),
	}

	k := keeper.NewKeeper(
		encCfg.Codec,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		addressCodec,
		mockKeepers.AccountKeeper,
		mockKeepers.BankKeeper,
		mockKeepers.StakingKeeper,
		mockKeepers.TokenConverterKeeper,
		mockKeepers.ShareclassKeeper,
	)

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
		mocks:        mockKeepers,
	}
}
