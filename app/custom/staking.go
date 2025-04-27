package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	staking "github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type CustomStakingModule struct {
	staking.AppModuleBasic
	cdc codec.Codec
}

func (cm CustomStakingModule) DefaultGenesis() json.RawMessage {
	genesis := stakingtypes.DefaultGenesisState()

	return cm.cdc.MustMarshalJSON(genesis)
}
