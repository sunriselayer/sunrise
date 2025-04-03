package shareclass

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	stakingtypes "cosmossdk.io/x/staking/types"

	"github.com/sunriselayer/sunrise/x/shareclass/keeper"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
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

	AccountKeeper        types.AccountKeeper
	BankKeeper           types.BankKeeper
	StakingKeeper        types.StakingKeeper
	FeeKeeper            types.FeeKeeper
	TokenConverterKeeper types.TokenConverterKeeper

	StakingMsgServer stakingtypes.MsgServer
}

type ModuleOutputs struct {
	depinject.Out

	ShareclassKeeper keeper.Keeper
	Module           appmodule.AppModule
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
		in.StakingKeeper,
		in.FeeKeeper,
		in.TokenConverterKeeper,
		in.StakingMsgServer,
	)
	m := NewAppModule(in.Cdc, k)

	return ModuleOutputs{ShareclassKeeper: k, Module: m}
}
