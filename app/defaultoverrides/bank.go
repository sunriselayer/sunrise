package defaultoverrides

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// BankModuleBasic defines a custom wrapper around the x/bank module's AppModuleBasic
// implementation to provide custom default genesis state.
type BankModuleBasic struct {
	bank.AppModuleBasic
}

// DefaultGenesis returns custom x/bank module genesis state.
func (BankModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	metadataFee := banktypes.Metadata{
		Description: "The native token of the Sunrise network for fees.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "urise",
				Exponent: 0,
				Aliases:  []string{"microrise"},
			},
			{
				Denom:    "rise",
				Exponent: 6,
			},
		},
		Base:    "urise",
		Display: "rise",
		Name:    "Sunrise RISE",
		Symbol:  "RISE",
	}
	metadataBond := banktypes.Metadata{
		Description: "The native token of the Sunrise network for staking. This token is non transferrable. This token can be retrieved by providing liquidity.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "uvrise",
				Exponent: 0,
				Aliases:  []string{"microvrise"},
			},
			{
				Denom:    "vrise",
				Exponent: 6,
			},
		},
		Base:    "uvrise",
		Display: "vrise",
		Name:    "Sunrise vRISE",
		Symbol:  "vRISE",
	}

	sendEnabledVrise := banktypes.SendEnabled{
		Denom:   "uvrise",
		Enabled: false,
	}

	genState := banktypes.DefaultGenesisState()
	genState.DenomMetadata = append(genState.DenomMetadata, metadataFee)
	genState.DenomMetadata = append(genState.DenomMetadata, metadataBond)

	genState.SendEnabled = append(genState.SendEnabled, sendEnabledVrise)

	return cdc.MustMarshalJSON(genState)
}
