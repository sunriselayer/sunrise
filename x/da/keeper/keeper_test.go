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

	"github.com/sunriselayer/sunrise/x/da/keeper"
	module "github.com/sunriselayer/sunrise/x/da/module"
	datestutil "github.com/sunriselayer/sunrise/x/da/testutil"
	"github.com/sunriselayer/sunrise/x/da/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("sunrise", "sunrisepub")
	config.SetBech32PrefixForValidator("sunrisevaloper", "sunrisevaloperpub")
	config.SetBech32PrefixForConsensusNode("sunrisevalcons", "sunrisevalconspub")

	encCfg := moduletestutil.MakeTestEncodingConfig(codectestutil.CodecOptions{}, module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(config.GetBech32AccountAddrPrefix())
	validatorAddressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	env := runtime.NewEnvironment(runtime.NewKVStoreService(storeKey), log.NewTestLogger(t))
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)

	mocks := getMocks(t)

	k := keeper.NewKeeper(
		env,
		encCfg.Codec,
		addressCodec,
		validatorAddressCodec,
		authority,
		mocks.BankKeeper,
		mocks.StakingKeeper,
		mocks.SlashingKeeper,
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

type DaMocks struct {
	BankKeeper     *datestutil.MockBankKeeper
	StakingKeeper  *datestutil.MockStakingKeeper
	SlashingKeeper *datestutil.MockSlashingKeeper
}

func getMocks(t *testing.T) DaMocks {

	ctrl := gomock.NewController(t)
	return DaMocks{
		BankKeeper:     datestutil.NewMockBankKeeper(ctrl),
		StakingKeeper:  datestutil.NewMockStakingKeeper(ctrl),
		SlashingKeeper: datestutil.NewMockSlashingKeeper(ctrl),
	}
}
