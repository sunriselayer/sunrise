package custombank

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	customkeeper "github.com/sunriselayer/sunrise/x/bank/keeper"
)

type AppModule struct {
	bank.AppModule
	customKeeper customkeeper.Keeper
}

func NewAppModule(bankModule bank.AppModule, customKeeper customkeeper.Keeper) AppModule {
	return AppModule{
		AppModule:    bankModule,
		customKeeper: customKeeper,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.customKeeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.customKeeper)
}
