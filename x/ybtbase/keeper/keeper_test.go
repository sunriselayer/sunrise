package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/golang/mock/gomock"

	"github.com/sunriselayer/sunrise/x/ybtbase/keeper"
	module "github.com/sunriselayer/sunrise/x/ybtbase/module"
	testtool "github.com/sunriselayer/sunrise/x/ybtbase/testutil"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

var (
	testAddress  = "cosmos1w3jhxap3gempvr"
	testAddress2 = "cosmos1w3jhxapjx2whzu"
	testAddress3 = "cosmos1w3jhxapnmu6zlw"
)

type moduleMocks struct {
	BankKeeper *testtool.MockBankKeeper
	AuthKeeper *testtool.MockAuthKeeper
}

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	mocks        moduleMocks
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)

	// Create mocks
	mocks := moduleMocks{
		BankKeeper: testtool.NewMockBankKeeper(ctrl),
		AuthKeeper: testtool.NewMockAuthKeeper(ctrl),
	}

	k := keeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		authority,
		mocks.AuthKeeper,
		mocks.BankKeeper,
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
