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

	"github.com/sunriselayer/sunrise/x/shareclass/keeper"
	module "github.com/sunriselayer/sunrise/x/shareclass/module"
	shareclasstestutil "github.com/sunriselayer/sunrise/x/shareclass/testutil"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	mocks        ShareclassMocks
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
		mocks.StakingKeeper,
		mocks.DistributionKeeper,
		mocks.TokenConverterKeeper,
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
		mocks:        mocks,
	}
}

type ShareclassMocks struct {
	AccountKeeper        *shareclasstestutil.MockAccountKeeper
	BankKeeper           *shareclasstestutil.MockBankKeeper
	DistributionKeeper   *shareclasstestutil.MockDistributionKeeper
	StakingKeeper        *shareclasstestutil.MockStakingKeeper
	TokenConverterKeeper *shareclasstestutil.MockTokenConverterKeeper
}

func getMocks(t *testing.T) ShareclassMocks {
	ctrl := gomock.NewController(t)

	return ShareclassMocks{
		AccountKeeper:        shareclasstestutil.NewMockAccountKeeper(ctrl),
		BankKeeper:           shareclasstestutil.NewMockBankKeeper(ctrl),
		DistributionKeeper:   shareclasstestutil.NewMockDistributionKeeper(ctrl),
		StakingKeeper:        shareclasstestutil.NewMockStakingKeeper(ctrl),
		TokenConverterKeeper: shareclasstestutil.NewMockTokenConverterKeeper(ctrl),
	}
}
