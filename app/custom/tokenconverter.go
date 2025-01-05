package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	tokenconverter "github.com/sunriselayer/sunrise/x/tokenconverter/module"
	tokenconvertertypes "github.com/sunriselayer/sunrise/x/tokenconverter/types"

	"github.com/sunriselayer/sunrise/app/consts"
)

type CustomTokenConverterModule struct {
	tokenconverter.AppModule
	cdc codec.Codec
}

func (cm CustomTokenConverterModule) DefaultGenesis() json.RawMessage {
	genesis := tokenconvertertypes.DefaultGenesis()

	genesis.Params.BondDenom = consts.BondDenom
	genesis.Params.FeeDenom = consts.FeeDenom

	return cm.cdc.MustMarshalJSON(genesis)
}
