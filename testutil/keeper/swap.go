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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/swap/keeper"
	swaptestutil "github.com/sunriselayer/sunrise/x/swap/testutil"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

type SwapMocks struct {
	AcctKeeper          *swaptestutil.MockAccountKeeper
	BankKeeper          *swaptestutil.MockBankKeeper
	TransferKeeper      *swaptestutil.MockTransferKeeper
	LiquiditypoolKeeper *swaptestutil.MockLiquidityPoolKeeper
}

func SwapKeeper(t testing.TB) (keeper.Keeper, SwapMocks, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	ctrl := gomock.NewController(t)
	m := SwapMocks{
		AcctKeeper:          swaptestutil.NewMockAccountKeeper(ctrl),
		BankKeeper:          swaptestutil.NewMockBankKeeper(ctrl),
		TransferKeeper:      swaptestutil.NewMockTransferKeeper(ctrl),
		LiquiditypoolKeeper: swaptestutil.NewMockLiquidityPoolKeeper(ctrl),
	}

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		m.AcctKeeper,
		m.BankKeeper,
		m.TransferKeeper,
		m.LiquiditypoolKeeper,
		nil,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, m, ctx
}
