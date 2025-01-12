package custom

import (
	"encoding/json"

	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/codec"

	staking "cosmossdk.io/x/staking"
	stakingkeeper "cosmossdk.io/x/staking/keeper"
	stakingtypes "cosmossdk.io/x/staking/types"

	customkeeper "github.com/sunriselayer/sunrise/app/custom/staking/keeper"
	customtypes "github.com/sunriselayer/sunrise/app/custom/types"
)

type CustomStakingModule struct {
	staking.AppModule
	cdc                  codec.Codec
	keeper               customtypes.StakingKeeper
	tokenconverterKeeper customtypes.TokenConverterKeeper
}

func (cm CustomStakingModule) DefaultGenesis() json.RawMessage {
	genesis := stakingtypes.DefaultGenesisState()

	genesis.Params.KeyRotationFee.Denom = genesis.Params.BondDenom

	return cm.cdc.MustMarshalJSON(genesis)
}

// RegisterServices registers module services.
func (cm CustomStakingModule) RegisterServices(registrar grpc.ServiceRegistrar) error {
	stakingtypes.RegisterMsgServer(registrar, customkeeper.NewMsgServerImpl(cm.keeper, cm.tokenconverterKeeper))
	stakingtypes.RegisterQueryServer(registrar, stakingkeeper.NewQuerier(cm.keeper.(*stakingkeeper.Keeper)))

	return nil
}
