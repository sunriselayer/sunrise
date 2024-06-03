package testencoding

import (
	sdkmodule "github.com/cosmos/cosmos-sdk/types/module"

	"github.com/sunriselayer/sunrise/app"
	"github.com/sunriselayer/sunrise/app/encoding"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
)

var (
	ModuleBasics = getModuleBasics()

	// ModuleEncodingRegisters keeps track of all the module methods needed to
	// register interfaces and specific type to encoding config
	ModuleEncodingRegisters = extractRegisters(ModuleBasics)
)

func getModuleBasics() sdkmodule.BasicManager {
	var moduleBasics sdkmodule.BasicManager
	depinject.Inject(
		depinject.Configs(
			app.AppConfig(),
			depinject.Supply(
				log.NewNopLogger(),
			)),
		&moduleBasics,
	)

	return moduleBasics
}

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

// // GetTxConfig implements the TestingApp interface.
// func (app *App) GetTxConfig() client.TxConfig {
// 	return app.txConfig
// }
