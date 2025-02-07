package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
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
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(codectestutil.CodecOptions{}, module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	env := runtime.NewEnvironment(runtime.NewKVStoreService(storeKey), log.NewTestLogger(t))
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)

	k := keeper.NewKeeper(
		env,
		encCfg.Codec,
		addressCodec,
		authority,
		nil,
		nil,
		nil,
		nil,
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
