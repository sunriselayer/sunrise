package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"cosmossdk.io/x/staking"
	stakingtypes "cosmossdk.io/x/staking/types"
)

type CustomStakingModule struct {
	staking.AppModule
	cdc codec.Codec
}

func (cm CustomStakingModule) DefaultGenesis() json.RawMessage {
	genesis := stakingtypes.DefaultGenesisState()

	genesis.Params.KeyRotationFee.Denom = genesis.Params.BondDenom

	return cm.cdc.MustMarshalJSON(genesis)
}
