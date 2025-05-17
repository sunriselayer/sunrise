package testutil

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

	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	module "github.com/sunriselayer/sunrise/x/liquidityincentive/module"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

type Fixture struct {
	Ctx          context.Context
	Keeper       keeper.Keeper
	AddressCodec address.Codec
	Mocks        LiquidityIncentiveMocks
}

func InitFixture(t *testing.T) *Fixture {
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
		mocks.StakingKeeper,
		mocks.FeeKeeper,
		mocks.TokenConverterKeeper,
		mocks.LiquidityPoolKeeper,
	)

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &Fixture{
		Ctx:          ctx,
		Keeper:       k,
		AddressCodec: addressCodec,
		Mocks:        mocks,
	}
}

type LiquidityIncentiveMocks struct {
	AccountKeeper        *MockAccountKeeper
	BankKeeper           *MockBankKeeper
	StakingKeeper        *MockStakingKeeper
	FeeKeeper            *MockFeeKeeper
	TokenConverterKeeper *MockTokenConverterKeeper
	LiquidityPoolKeeper  *MockLiquidityPoolKeeper
}

func getMocks(t *testing.T) LiquidityIncentiveMocks {
	ctrl := gomock.NewController(t)
	return LiquidityIncentiveMocks{
		AccountKeeper:        NewMockAccountKeeper(ctrl),
		BankKeeper:           NewMockBankKeeper(ctrl),
		StakingKeeper:        NewMockStakingKeeper(ctrl),
		FeeKeeper:            NewMockFeeKeeper(ctrl),
		TokenConverterKeeper: NewMockTokenConverterKeeper(ctrl),
		LiquidityPoolKeeper:  NewMockLiquidityPoolKeeper(ctrl),
	}
}
