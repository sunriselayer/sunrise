package staking

import (
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/codec"

	stakingkeeper "cosmossdk.io/x/staking/keeper"
	stakingtypes "cosmossdk.io/x/staking/types"

	tokenconverterkeeper "github.com/sunriselayer/sunrise/x/tokenconverter/keeper"

	customkeeper "github.com/sunriselayer/sunrise/x/staking/keeper"
)

type CustomStakingModule struct {
	cdc                  codec.Codec
	keeper               *stakingkeeper.Keeper
	tokenconverterKeeper tokenconverterkeeper.Keeper
}

// RegisterServices registers module services.
func (cm CustomStakingModule) RegisterServices(registrar grpc.ServiceRegistrar) error {
	stakingtypes.RegisterMsgServer(registrar, customkeeper.NewMsgServerImpl(cm.keeper, cm.tokenconverterKeeper))
	stakingtypes.RegisterQueryServer(registrar, stakingkeeper.NewQuerier(cm.keeper))

	return nil
}
