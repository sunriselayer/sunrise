package ybtbrand

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/sunriselayer/sunrise/x/ybtbrand/keeper"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec
	Logger       log.Logger

	AuthKeeper    types.AuthKeeper
	BankKeeper    types.BankKeeper
	YbtbaseKeeper types.YbtbaseKeeper
}

type ModuleOutputs struct {
	depinject.Out

	YbtbrandKeeper keeper.Keeper
	Module         appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k, err := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Logger,
		in.AuthKeeper,
		in.BankKeeper,
		in.YbtbaseKeeper,
		in.AddressCodec,
	)
	if err != nil {
		panic(err)
	}

	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{YbtbrandKeeper: k, Module: m}
}
