package custombank

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"

	modulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	"github.com/cosmos/cosmos-sdk/x/bank"

	customkeeper "github.com/sunriselayer/sunrise/x/bank/keeper"
)

func init() {
	appconfig.Register(
		&modulev1.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleOutputs struct {
	depinject.Out

	BankKeeper customkeeper.Keeper
	Module     appmodule.AppModule
}

func ProvideModule(in bank.ModuleInputs) ModuleOutputs {
	outputs := bank.ProvideModule(in)

	k := customkeeper.NewKeeper(
		outputs.BankKeeper,
		in.AccountKeeper,
		nil,
	)
	bankModule, ok := outputs.Module.(bank.AppModule)
	if !ok {
		panic("bank module is not a bank.AppModule")
	}
	m := NewAppModule(bankModule, k)

	return ModuleOutputs{
		BankKeeper: k,
		Module:     m,
	}
}
