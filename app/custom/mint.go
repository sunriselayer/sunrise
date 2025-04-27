package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/sunriselayer/sunrise/app/consts"
)

type CustomMintModule struct {
	mint.AppModuleBasic
	cdc codec.Codec
}

func (cm CustomMintModule) DefaultGenesis() json.RawMessage {
	genesis := minttypes.DefaultGenesisState()

	// Params wil not be used anyway because there is a custom MintFn
	genesis.Params.MintDenom = consts.FeeDenom

	return cm.cdc.MustMarshalJSON(genesis)
}
