package custom

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/x/protocolpool"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"

	"github.com/sunriselayer/sunrise/app/consts"
)

type CustomProtocolPoolModule struct {
	protocolpool.AppModule
}

// DefaultGenesis returns default genesis state as raw bytes for the protocolpool module.
func (cm CustomProtocolPoolModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := protocolpooltypes.DefaultGenesisState()

	// Params wil not be used anyway because there is a custom MintFn
	genesis.Params.EnabledDistributionDenoms = []string{consts.NativeDenom, consts.StableDenom}

	return cdc.MustMarshalJSON(genesis)
}

// ValidateGenesis performs genesis state validation for the protocolpool module.
func (cm CustomProtocolPoolModule) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState protocolpooltypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", protocolpooltypes.ModuleName, err)
	}

	return genState.Validate()
}
