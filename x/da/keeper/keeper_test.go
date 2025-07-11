package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/da/keeper"
	module "github.com/sunriselayer/sunrise/x/da/module"
	datestutil "github.com/sunriselayer/sunrise/x/da/testutil"
	"github.com/sunriselayer/sunrise/x/da/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	mocks        DaMocks
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	config := sdk.GetConfig()
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(config.GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)

	mocks := getMocks(t)

	k := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		log.NewNopLogger(),
		authority.String(),
		addressCodec,
		mocks.BankKeeper,
		mocks.StakingKeeper,
		mocks.SlashingKeeper,
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

type DaMocks struct {
	BankKeeper     *datestutil.MockBankKeeper
	StakingKeeper  *datestutil.MockStakingKeeper
	SlashingKeeper *datestutil.MockSlashingKeeper
}

func getMocks(t *testing.T) DaMocks {

	ctrl := gomock.NewController(t)
	return DaMocks{
		BankKeeper:     datestutil.NewMockBankKeeper(ctrl),
		StakingKeeper:  datestutil.NewMockStakingKeeper(ctrl),
		SlashingKeeper: datestutil.NewMockSlashingKeeper(ctrl),
	}
}

// mockIterator helps in mocking the staking keeper's iterator
type mockIterator struct {
	valAddrs []sdk.ValAddress
	cursor   int
}

func (m *mockIterator) Domain() (start, end []byte) {
	return nil, nil
}

func (m *mockIterator) Valid() bool {
	return m.cursor < len(m.valAddrs)
}

func (m *mockIterator) Next() {
	m.cursor++
}

func (m *mockIterator) Key() (key []byte) {
	return []byte{}
}

func (m *mockIterator) Value() (value []byte) {
	if m.cursor >= len(m.valAddrs) {
		return nil
	}
	return m.valAddrs[m.cursor]
}

func (m *mockIterator) Error() error {
	return nil
}

func (m *mockIterator) Close() error {
	return nil
}
