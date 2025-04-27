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

	"github.com/sunriselayer/sunrise/x/swap/keeper"
	module "github.com/sunriselayer/sunrise/x/swap/module"
	swaptestutil "github.com/sunriselayer/sunrise/x/swap/testutil"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	mocks        SwapMocks
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
		mocks.TransferKeeper,
		mocks.LiquiditypoolKeeper,
		nil,
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

type SwapMocks struct {
	AcctKeeper          *swaptestutil.MockAccountKeeper
	BankKeeper          *swaptestutil.MockBankKeeper
	TransferKeeper      *swaptestutil.MockTransferKeeper
	LiquiditypoolKeeper *swaptestutil.MockLiquidityPoolKeeper
}

func getMocks(t *testing.T) SwapMocks {

	ctrl := gomock.NewController(t)
	return SwapMocks{
		AcctKeeper:          swaptestutil.NewMockAccountKeeper(ctrl),
		BankKeeper:          swaptestutil.NewMockBankKeeper(ctrl),
		TransferKeeper:      swaptestutil.NewMockTransferKeeper(ctrl),
		LiquiditypoolKeeper: swaptestutil.NewMockLiquidityPoolKeeper(ctrl),
	}
}
