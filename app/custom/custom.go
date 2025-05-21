package custom

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/protocolpool"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fee "github.com/sunriselayer/sunrise/x/fee/module"
	feetypes "github.com/sunriselayer/sunrise/x/fee/types"

	"github.com/sunriselayer/sunrise/app/consts"
)

func ReplaceCustomModules(
	manager module.BasicManager,
) {
	sdk.DefaultBondDenom = consts.BondDenom

	// bank
	oldBankModule, _ := manager[banktypes.ModuleName].(bank.AppModuleBasic)
	manager[banktypes.ModuleName] = CustomBankModule{
		AppModuleBasic: oldBankModule,
	}

	// fee
	oldFeeModule, _ := manager[feetypes.ModuleName].(fee.AppModuleBasic)
	manager[feetypes.ModuleName] = CustomFeeModule{
		AppModuleBasic: oldFeeModule,
	}

	// gov
	oldGovModule, _ := manager[govtypes.ModuleName].(gov.AppModuleBasic)
	manager[govtypes.ModuleName] = CustomGovModule{
		AppModuleBasic: oldGovModule,
	}

	// mint
	oldMintModule, _ := manager[minttypes.ModuleName].(mint.AppModuleBasic)
	manager[minttypes.ModuleName] = CustomMintModule{
		AppModuleBasic: oldMintModule,
	}

	// protocolpool
	oldProtocolPoolModule, _ := manager[protocolpooltypes.ModuleName].(protocolpool.AppModule)
	manager[protocolpooltypes.ModuleName] = CustomProtocolPoolModule{
		AppModule: oldProtocolPoolModule,
	}

	// staking
	oldStakingModule, _ := manager[stakingtypes.ModuleName].(staking.AppModuleBasic)
	manager[stakingtypes.ModuleName] = CustomStakingModule{
		AppModuleBasic: oldStakingModule,
	}
}
