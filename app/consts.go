package app

import (
	sdkmodule "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/sunrise-zone/sunrise-app/app/encoding"
	"github.com/sunrise-zone/sunrise-app/pkg/appconsts"

	"github.com/cosmos/cosmos-sdk/client"
)

const (
	// BondDenom defines the native staking token denomination.
	BondDenom = appconsts.BondDenom
	// BondDenomAlias defines an alias for BondDenom.
	BondDenomAlias = "micro-sr"
	// DisplayDenom defines the name, symbol, and display value of the Celestia token.
	DisplayDenom = "SR"
)

var (
	ModuleBasics = sdkmodule.NewBasicManager()
	// ModuleEncodingRegisters keeps track of all the module methods needed to
	// register interfaces and specific type to encoding config
	ModuleEncodingRegisters = extractRegisters(ModuleBasics)
)

// extractRegisters isolates the encoding module registers from the module
// manager, and appends any solo registers.
func extractRegisters(m sdkmodule.BasicManager, soloRegisters ...encoding.ModuleRegister) []encoding.ModuleRegister {
	// TODO: might be able to use some standard generics in go 1.18
	s := make([]encoding.ModuleRegister, len(m)+len(soloRegisters))
	i := 0
	for _, v := range m {
		s[i] = v
		i++
	}
	for i, v := range soloRegisters {
		s[i+len(m)] = v
	}
	return s
}

// GetTxConfig implements the TestingApp interface.
func (app *App) GetTxConfig() client.TxConfig {
	return app.txConfig
}
