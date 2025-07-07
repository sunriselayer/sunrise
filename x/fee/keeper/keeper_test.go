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

	"github.com/sunriselayer/sunrise/x/fee/keeper"
	module "github.com/sunriselayer/sunrise/x/fee/module"
	feetestutil "github.com/sunriselayer/sunrise/x/fee/testutil"
	"github.com/sunriselayer/sunrise/x/fee/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	mocks        feeMocks
}

type feeMocks struct {
	AccountKeeper       *feetestutil.MockAccountKeeper
	BankKeeper          *feetestutil.MockBankKeeper
	LiquiditypoolKeeper *feetestutil.MockLiquidityPoolKeeper
	SwapKeeper          *feetestutil.MockSwapKeeper
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
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
		mocks.AccountKeeper,
		mocks.BankKeeper,
		mocks.LiquiditypoolKeeper,
		mocks.SwapKeeper,
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

func getMocks(t *testing.T) feeMocks {
	ctrl := gomock.NewController(t)

	return feeMocks{
		AccountKeeper:       feetestutil.NewMockAccountKeeper(ctrl),
		BankKeeper:          feetestutil.NewMockBankKeeper(ctrl),
		LiquiditypoolKeeper: feetestutil.NewMockLiquidityPoolKeeper(ctrl),
		SwapKeeper:          feetestutil.NewMockSwapKeeper(ctrl),
	}
}
