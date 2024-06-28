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
	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	liquiditypooltestutil "github.com/sunriselayer/sunrise/x/liquiditypool/testutil"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("sunrise", "sunrisepub")
	config.SetBech32PrefixForValidator("sunrisevaloper", "sunrisevaloperpub")
}

func LiquiditypoolKeeper(t *testing.T) (keeper.Keeper, *liquiditypooltestutil.MockBankKeeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	ctrl := gomock.NewController(t)
	bankKeeper := liquiditypooltestutil.NewMockBankKeeper(ctrl)

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		bankKeeper,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, bankKeeper, ctx
}
