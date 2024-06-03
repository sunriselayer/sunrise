package defaultoverrides

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// MintModuleBasic defines a custom wrapper around the x/bank module's AppModuleBasic
// implementation to provide custom default genesis state.
type MintModuleBasic struct {
	mint.AppModuleBasic
}

// DefaultGenesis returns custom x/bank module genesis state.
func (MintModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genState := minttypes.DefaultGenesisState()
	genState.Params.MintDenom = "uvrise"

	return cdc.MustMarshalJSON(genState)
}
