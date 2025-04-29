package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	fee "github.com/sunriselayer/sunrise/x/fee/module"
	feetypes "github.com/sunriselayer/sunrise/x/fee/types"

	"github.com/sunriselayer/sunrise/app/consts"
)

type CustomFeeModule struct {
	fee.AppModuleBasic
}

func (cm CustomFeeModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := feetypes.DefaultGenesis()

	genesis.Params.FeeDenom = consts.FeeDenom

	return cdc.MustMarshalJSON(genesis)
}
