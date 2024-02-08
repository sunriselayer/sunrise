package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmodule "github.com/cosmos/cosmos-sdk/types/module"

	"github.com/sunrise-zone/sunrise-app/app/encoding"
	"github.com/sunrise-zone/sunrise-app/pkg/appconsts"
	blobmodule "github.com/sunrise-zone/sunrise-app/x/blob/module"
	bsmodule "github.com/sunrise-zone/sunrise-app/x/blobstream/module"

	// "cosmossdk.io/depinject"
	"cosmossdk.io/x/evidence"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"

	"github.com/cosmos/cosmos-sdk/client"
	auth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	mint "github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"

	capability "github.com/cosmos/ibc-go/modules/capability"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
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
	// ModuleEncodingRegisters keeps track of all the module methods needed to
	// register interfaces and specific type to encoding config
	ModuleEncodingRegisters = extractRegisters(ModuleBasics())
)

func ModuleBasics() sdkmodule.BasicManager {
	// Set prefixes
	accountPubKeyPrefix := Bech32PrefixAccPub
	validatorAddressPrefix := Bech32PrefixValAddr
	validatorPubKeyPrefix := Bech32PrefixValPub
	consNodeAddressPrefix := Bech32PrefixConsAddr
	consNodePubKeyPrefix := Bech32PrefixConsPub

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)

	moduleBasics := sdkmodule.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bankModule{},
		capability.AppModuleBasic{},
		stakingModule{},
		mint.AppModuleBasic{},
		distributionModule{},
		govModule{},
		params.AppModuleBasic{},
		crisisModule{},
		slashingModule{},
		authzmodule.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		blobmodule.AppModuleBasic{},
		bsmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
	)
	// depinject.Inject(
	// 	depinject.Configs(AppConfig()),
	// 	&moduleBasics,
	// )

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

// GetTxConfig implements the TestingApp interface.
func (app *App) GetTxConfig() client.TxConfig {
	return app.txConfig
}
