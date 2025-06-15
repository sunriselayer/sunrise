package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	ybtbrand "github.com/sunriselayer/sunrise/x/ybtbrand/module"
	"github.com/sunriselayer/sunrise/x/ybtbrand/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbrand/testutil"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

const (
	testAddress  = "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu"
	testAddress2 = "cosmos1z7g5w84ynmjyg0kqpahdjqpj7yq34v3suckp0e"
	testAddress3 = "cosmos1s0he0z3g92zwsxdj83h0ky9w463sx7gq9mqtgn"
)

type fixture struct {
	ctx         context.Context
	keeper      keeper.Keeper
	queryServer types.QueryServer
	msgServer   types.MsgServer
	mocks       moduleMocks
}

type moduleMocks struct {
	BankKeeper    *testutil.MockBankKeeper
	AuthKeeper    *testutil.MockAuthKeeper
	YbtbaseKeeper *testutil.MockYbtbaseKeeper
}

func initFixture(t *testing.T) *fixture {
	ctrl := gomock.NewController(t)
	mockBankKeeper := testutil.NewMockBankKeeper(ctrl)
	mockAuthKeeper := testutil.NewMockAuthKeeper(ctrl)
	mockYbtbaseKeeper := testutil.NewMockYbtbaseKeeper(ctrl)

	mocks := moduleMocks{
		BankKeeper:    mockBankKeeper,
		AuthKeeper:    mockAuthKeeper,
		YbtbaseKeeper: mockYbtbaseKeeper,
	}

	keys := storetypes.NewKVStoreKeys(types.StoreKey)
	cdc := moduletestutil.MakeTestEncodingConfig(ybtbrand.AppModule{}).Codec

	logger := log.NewTestLogger(t)
	cms := integration.CreateMultiStore(keys, logger)

	newCtx := sdk.NewContext(cms, cmtproto.Header{}, true, logger)

	k, err := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keys[types.StoreKey]),
		logger,
		mockAuthKeeper,
		mockBankKeeper,
		mockYbtbaseKeeper,
		address.NewBech32Codec("cosmos"),
	)
	require.NoError(t, err)

	ctx := newCtx

	return &fixture{
		ctx:         ctx,
		keeper:      k,
		queryServer: keeper.NewQueryServerImpl(k),
		msgServer:   keeper.NewMsgServerImpl(k),
		mocks:       mocks,
	}
}

// CreateMultiStore creates a multi-store for testing
func CreateMultiStore(keys map[string]*storetypes.KVStoreKey, logger log.Logger) storetypes.CommitMultiStore {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db, logger, metrics.NewNoOpMetrics())
	for key, storeKey := range keys {
		cms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
		logger.Info("Mounted store", "key", key)
	}
	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	return cms
}
