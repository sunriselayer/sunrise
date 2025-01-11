package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"cosmossdk.io/x/protocolpool"
	protocolpooltypes "cosmossdk.io/x/protocolpool/types"

	"github.com/sunriselayer/sunrise/app/consts"
)

type CustomProtocolPoolModule struct {
	protocolpool.AppModule
	cdc codec.Codec
}

func (cm CustomProtocolPoolModule) DefaultGenesis() json.RawMessage {
	genesis := protocolpooltypes.DefaultGenesisState()

	// Params wil not be used anyway because there is a custom MintFn
	genesis.Params.EnabledDistributionDenoms = []string{consts.FeeDenom}

	return cm.cdc.MustMarshalJSON(genesis)
}
