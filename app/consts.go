package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmodule "github.com/cosmos/cosmos-sdk/types/module"

	"github.com/sunriselayer/sunrise/app/encoding"
	blobmodule "github.com/sunriselayer/sunrise/x/blob/module"
	bsmodule "github.com/sunriselayer/sunrise/x/blobstream/module"
	feemodule "github.com/sunriselayer/sunrise/x/fee/module"
	liquidityincentivemodule "github.com/sunriselayer/sunrise/x/liquidityincentive/module"
	liquiditypoolmodule "github.com/sunriselayer/sunrise/x/liquiditypool/module"
	swapmodule "github.com/sunriselayer/sunrise/x/swap/module"
	tokenconvertermodule "github.com/sunriselayer/sunrise/x/tokenconverter/module"

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

var (
	// ModuleEncodingRegisters keeps track of all the module methods needed to
	// register interfaces and specific type to encoding config
	ModuleEncodingRegisters = extractRegisters(ModuleBasics())
)

// These constants are derived from the above variables.
// These are the ones we will want to use in the code, based on
// any overrides above.
var (
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address.
	Bech32PrefixAccAddr = AccountAddressPrefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key.
	Bech32PrefixAccPub = AccountAddressPrefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address.
	Bech32PrefixValAddr = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key.
	Bech32PrefixValPub = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address.
	Bech32PrefixConsAddr = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key.
	Bech32PrefixConsPub = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
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
		tokenconvertermodule.AppModuleBasic{},
		liquiditypoolmodule.AppModuleBasic{},
		liquidityincentivemodule.AppModuleBasic{},
		swapmodule.AppModuleBasic{},
		feemodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
	)
	// moduleBasics := sdkmodule.BasicManager{}
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
