package swap

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/sunriselayer/sunrise/x/swap/keeper"
	"github.com/sunriselayer/sunrise/x/swap/types"

	ibckeeper "github.com/cosmos/ibc-go/v9/modules/core/keeper"
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
	Environment  appmodule.Environment
	Cdc          codec.Codec
	AddressCodec address.Codec

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
	// TransferKeeper      types.TransferKeeper
	LiquidityPoolKeeper types.LiquidityPoolKeeper

	IBCKeeperFn func() *ibckeeper.Keeper `optional:"true"`
}

type ModuleOutputs struct {
	depinject.Out

	SwapKeeper keeper.Keeper
	Module     appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(types.GovModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.Environment,
		in.Cdc,
		in.AddressCodec,
		authority,
		in.AccountKeeper,
		in.BankKeeper,
		nil, // TransferKeeper is not supported by DepInject
		in.LiquidityPoolKeeper,
		in.IBCKeeperFn,
	)
	m := NewAppModule(in.Cdc, k)

	return ModuleOutputs{SwapKeeper: k, Module: m}
}
