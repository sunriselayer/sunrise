package defaultoverrides

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	feemodule "github.com/sunriselayer/sunrise/x/fee/module"
	feetypes "github.com/sunriselayer/sunrise/x/fee/types"
)

// FeeModuleBasic defines a custom wrapper around the x/fee module's AppModuleBasic
// implementation to provide custom default genesis state.
type FeeModuleBasic struct {
	feemodule.AppModuleBasic
}

// DefaultGenesis returns custom x/fee module genesis state.
func (FeeModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genState := feetypes.DefaultGenesis()
	genState.Params.FeeDenom = "urise"
	genState.Params.BypassDenoms = []string{"uvrise"}

	return cdc.MustMarshalJSON(genState)
}
