package defaultoverrides

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	feetypes "github.com/sunriselayer/sunrise/x/fee/types"
)

// FeeModuleBasic defines a custom wrapper around the x/fee module's AppModuleBasic
// implementation to provide custom default genesis state.
type FeeModuleBasic struct {
	bank.AppModuleBasic
}

// DefaultGenesis returns custom x/fee module genesis state.
func (FeeModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := feetypes.DefaultGenesis()
	genesis.Params.FeeDenom = "urise"
	genesis.Params.BypassDenoms = []string{"uvrise"}

	return cdc.MustMarshalJSON(genesis)
}
