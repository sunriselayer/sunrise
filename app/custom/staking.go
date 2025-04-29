package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	staking "github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type CustomStakingModule struct {
	staking.AppModuleBasic
}

func (cm CustomStakingModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := stakingtypes.DefaultGenesisState()

	return cdc.MustMarshalJSON(genesis)
}
