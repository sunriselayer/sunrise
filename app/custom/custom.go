package custom

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"cosmossdk.io/x/bank"
	banktypes "cosmossdk.io/x/bank/types"
	"cosmossdk.io/x/gov"
	govtypes "cosmossdk.io/x/gov/types"
	"cosmossdk.io/x/mint"
	minttypes "cosmossdk.io/x/mint/types"
	"cosmossdk.io/x/protocolpool"
	protocolpooltypes "cosmossdk.io/x/protocolpool/types"

	"cosmossdk.io/x/staking"
	stakingtypes "cosmossdk.io/x/staking/types"

	fee "github.com/sunriselayer/sunrise/x/fee/module"
	feetypes "github.com/sunriselayer/sunrise/x/fee/types"
	tokenconverter "github.com/sunriselayer/sunrise/x/tokenconverter/module"
	tokenconvertertypes "github.com/sunriselayer/sunrise/x/tokenconverter/types"

	"github.com/sunriselayer/sunrise/app/consts"
)

func ReplaceCustomModules(
	manager *module.Manager,
	cdc codec.Codec,
) {
	sdk.DefaultBondDenom = consts.BondDenom

	// bank
	oldBankModule, _ := manager.Modules[banktypes.ModuleName].(bank.AppModule)
	manager.Modules[banktypes.ModuleName] = CustomBankModule{
		AppModule: oldBankModule,
		cdc:       cdc,
	}

	// fee
	oldFeeModule, _ := manager.Modules[feetypes.ModuleName].(fee.AppModule)
	manager.Modules[feetypes.ModuleName] = CustomFeeModule{
		AppModule: oldFeeModule,
		cdc:       cdc,
	}

	// gov
	oldGovModule, _ := manager.Modules[govtypes.ModuleName].(gov.AppModule)
	manager.Modules[govtypes.ModuleName] = CustomGovModule{
		AppModule: oldGovModule,
		cdc:       cdc,
	}

	// mint
	oldMintModule, _ := manager.Modules[minttypes.ModuleName].(mint.AppModule)
	manager.Modules[minttypes.ModuleName] = CustomMintModule{
		AppModule: oldMintModule,
		cdc:       cdc,
	}

	// protocolpool
	oldProtocolPoolModule, _ := manager.Modules[protocolpooltypes.ModuleName].(protocolpool.AppModule)
	manager.Modules[protocolpooltypes.ModuleName] = CustomProtocolPoolModule{
		AppModule: oldProtocolPoolModule,
		cdc:       cdc,
	}

	// staking
	oldStakingModule, _ := manager.Modules[stakingtypes.ModuleName].(staking.AppModule)
	manager.Modules[stakingtypes.ModuleName] = CustomStakingModule{
		AppModule: oldStakingModule,
		cdc:       cdc,
	}

	// tokenconverter
	oldTokenConverterModule, _ := manager.Modules[tokenconvertertypes.ModuleName].(tokenconverter.AppModule)
	manager.Modules[tokenconvertertypes.ModuleName] = CustomTokenConverterModule{
		AppModule: oldTokenConverterModule,
		cdc:       cdc,
	}
}
