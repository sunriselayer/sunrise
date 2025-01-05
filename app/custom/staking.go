package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"cosmossdk.io/x/staking"
	stakingtypes "cosmossdk.io/x/staking/types"

	"github.com/sunriselayer/sunrise/app/consts"
)

type CustomStakingModule struct {
	staking.AppModule
	cdc codec.Codec
}

func (cm CustomStakingModule) DefaultGenesis() json.RawMessage {
	genesis := stakingtypes.DefaultGenesisState()

	genesis.Params.BondDenom = consts.BondDenom

	return cm.cdc.MustMarshalJSON(genesis)
}
