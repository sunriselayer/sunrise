package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/x/protocolpool"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"

	"github.com/sunriselayer/sunrise/app/consts"
)

type CustomProtocolPoolModule struct {
	protocolpool.AppModule
}

func (cm CustomProtocolPoolModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := protocolpooltypes.DefaultGenesisState()

	// Params wil not be used anyway because there is a custom MintFn
	genesis.Params.EnabledDistributionDenoms = []string{consts.FeeDenom}

	return cdc.MustMarshalJSON(genesis)
}
