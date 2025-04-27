package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type CustomBankModule struct {
	bank.AppModuleBasic
	cdc codec.Codec
}

func (cm CustomBankModule) DefaultGenesis() json.RawMessage {
	genesis := banktypes.DefaultGenesisState()

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

	genesis.DenomMetadata = append(genesis.DenomMetadata, metadataFee)
	genesis.DenomMetadata = append(genesis.DenomMetadata, metadataBond)

	genesis.SendEnabled = append(genesis.SendEnabled, sendEnabledVrise)

	return cm.cdc.MustMarshalJSON(genesis)
}
