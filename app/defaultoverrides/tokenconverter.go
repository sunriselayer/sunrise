package defaultoverrides

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	tokenconvertertypes "github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

// TokenConverterModuleBasic defines a custom wrapper around the x/tokenconverter module's AppModuleBasic
// implementation to provide custom default genesis state.
type TokenConverterModuleBasic struct {
	bank.AppModuleBasic
}

// DefaultGenesis returns custom x/tokenconverter module genesis state.
func (m TokenConverterModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := tokenconvertertypes.DefaultGenesis()
	genesis.Params.BondDenom = "uvrise"
	genesis.Params.FeeDenom = "urise"

	return cdc.MustMarshalJSON(genesis)
}
