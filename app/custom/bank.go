package custom

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type CustomBankModule struct {
	bank.AppModuleBasic
}

func (cm CustomBankModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := banktypes.DefaultGenesisState()

	metadataMint := banktypes.Metadata{
		Description: "The native token of the Sunrise network.",
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
	metadataStable := banktypes.Metadata{
		Description: "The USD stable coin of the Sunrise network for fees.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "uusdrise",
				Exponent: 0,
				Aliases:  []string{"microusdrise"},
			},
			{
				Denom:    "usdrise",
				Exponent: 6,
			},
		},
		Base:    "uusdrise",
		Display: "usdrise",
		Name:    "Sunrise USDrise",
		Symbol:  "USDRISE",
	}

	sendEnabledMint := banktypes.SendEnabled{
		Denom:   "urise",
		Enabled: false,
	}
	sendEnabledVrise := banktypes.SendEnabled{
		Denom:   "uvrise",
		Enabled: false,
	}

	genesis.DenomMetadata = append(genesis.DenomMetadata, metadataMint)
	genesis.DenomMetadata = append(genesis.DenomMetadata, metadataBond)
	genesis.DenomMetadata = append(genesis.DenomMetadata, metadataStable)

	genesis.SendEnabled = append(genesis.SendEnabled, sendEnabledMint)
	genesis.SendEnabled = append(genesis.SendEnabled, sendEnabledVrise)

	return cdc.MustMarshalJSON(genesis)
}
