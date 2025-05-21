package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	module "github.com/sunriselayer/sunrise/x/liquidityincentive/module"
	liquidityincentivetestutil "github.com/sunriselayer/sunrise/x/liquidityincentive/testutil"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	mocks        LiquidityIncentiveMocks
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
		mocks.AcctKeeper,
		mocks.BankKeeper,
		mocks.StakingKeeper,
		mocks.FeeKeeper,
		mocks.TokenConverterKeeper,
		mocks.LiquiditypoolKeeper,
	)

	// Initialize params
	if err := k.Params.Set(ctx, types.NewParams(5, math.LegacyNewDecWithPrec(50, 2), 20)); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
		mocks:        mocks,
	}
}

type LiquidityIncentiveMocks struct {
	AcctKeeper           *liquidityincentivetestutil.MockAccountKeeper
	BankKeeper           *liquidityincentivetestutil.MockBankKeeper
	StakingKeeper        *liquidityincentivetestutil.MockStakingKeeper
	FeeKeeper            *liquidityincentivetestutil.MockFeeKeeper
	TokenConverterKeeper *liquidityincentivetestutil.MockTokenConverterKeeper
	LiquiditypoolKeeper  *liquidityincentivetestutil.MockLiquidityPoolKeeper
}

func getMocks(t *testing.T) LiquidityIncentiveMocks {
	ctrl := gomock.NewController(t)

	return LiquidityIncentiveMocks{
		AcctKeeper:           liquidityincentivetestutil.NewMockAccountKeeper(ctrl),
		BankKeeper:           liquidityincentivetestutil.NewMockBankKeeper(ctrl),
		StakingKeeper:        liquidityincentivetestutil.NewMockStakingKeeper(ctrl),
		FeeKeeper:            liquidityincentivetestutil.NewMockFeeKeeper(ctrl),
		TokenConverterKeeper: liquidityincentivetestutil.NewMockTokenConverterKeeper(ctrl),
		LiquiditypoolKeeper:  liquidityincentivetestutil.NewMockLiquidityPoolKeeper(ctrl),
	}
}
