package app

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	accountsmodulev1 "cosmossdk.io/api/cosmos/accounts/module/v1"
	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	authzmodulev1 "cosmossdk.io/api/cosmos/authz/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	circuitmodulev1 "cosmossdk.io/api/cosmos/circuit/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	distrmodulev1 "cosmossdk.io/api/cosmos/distribution/module/v1"
	epochsmodulev1 "cosmossdk.io/api/cosmos/epochs/module/v1"
	evidencemodulev1 "cosmossdk.io/api/cosmos/evidence/module/v1"
	feegrantmodulev1 "cosmossdk.io/api/cosmos/feegrant/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	govmodulev1 "cosmossdk.io/api/cosmos/gov/module/v1"
	groupmodulev1 "cosmossdk.io/api/cosmos/group/module/v1"
	mintmodulev1 "cosmossdk.io/api/cosmos/mint/module/v1"
	nftmodulev1 "cosmossdk.io/api/cosmos/nft/module/v1"
	paramsmodulev1 "cosmossdk.io/api/cosmos/params/module/v1"
	poolmodulev1 "cosmossdk.io/api/cosmos/protocolpool/module/v1"
	slashingmodulev1 "cosmossdk.io/api/cosmos/slashing/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	upgrademodulev1 "cosmossdk.io/api/cosmos/upgrade/module/v1"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/x/accounts"
	"cosmossdk.io/x/authz"
	_ "cosmossdk.io/x/authz/module" // import for side-effects
	_ "cosmossdk.io/x/bank"         // import for side-effects
	banktypes "cosmossdk.io/x/bank/types"
	_ "cosmossdk.io/x/circuit" // import for side-effects
	circuittypes "cosmossdk.io/x/circuit/types"
	_ "cosmossdk.io/x/consensus" // import for side-effects
	consensustypes "cosmossdk.io/x/consensus/types"
	_ "cosmossdk.io/x/distribution" // import for side-effects
	distrtypes "cosmossdk.io/x/distribution/types"
	_ "cosmossdk.io/x/epochs" // import for side-effects
	epochstypes "cosmossdk.io/x/epochs/types"
	_ "cosmossdk.io/x/evidence" // import for side-effects
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	_ "cosmossdk.io/x/feegrant/module" // import for side-effects
	_ "cosmossdk.io/x/gov"             // import for side-effects
	govtypes "cosmossdk.io/x/gov/types"
	"cosmossdk.io/x/group"
	_ "cosmossdk.io/x/group/module" // import for side-effects
	_ "cosmossdk.io/x/mint"         // import for side-effects
	minttypes "cosmossdk.io/x/mint/types"
	"cosmossdk.io/x/nft"
	_ "cosmossdk.io/x/nft/module" // import for side-effects
	_ "cosmossdk.io/x/params"     // import for side-effects
	paramstypes "cosmossdk.io/x/params/types"
	_ "cosmossdk.io/x/protocolpool" // import for side-effects
	pooltypes "cosmossdk.io/x/protocolpool/types"
	_ "cosmossdk.io/x/slashing" // import for side-effects
	slashingtypes "cosmossdk.io/x/slashing/types"
	_ "cosmossdk.io/x/staking" // import for side-effects
	stakingtypes "cosmossdk.io/x/staking/types"
	_ "cosmossdk.io/x/upgrade" // import for side-effects
	upgradetypes "cosmossdk.io/x/upgrade/types"

	_ "github.com/sunriselayer/sunrise/x/da/module"
	damoduletypes "github.com/sunriselayer/sunrise/x/da/types"
	_ "github.com/sunriselayer/sunrise/x/fee/module"
	feemoduletypes "github.com/sunriselayer/sunrise/x/fee/types"
	_ "github.com/sunriselayer/sunrise/x/liquidityincentive/module"
	liquidityincentivemoduletypes "github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	_ "github.com/sunriselayer/sunrise/x/liquiditypool/module"
	liquiditypoolmoduletypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	_ "github.com/sunriselayer/sunrise/x/selfdelegation/module"
	selfdelegationmoduletypes "github.com/sunriselayer/sunrise/x/selfdelegation/types"
	_ "github.com/sunriselayer/sunrise/x/swap/module"
	swapmoduletypes "github.com/sunriselayer/sunrise/x/swap/types"
	_ "github.com/sunriselayer/sunrise/x/tokenconverter/module"
	tokenconvertermoduletypes "github.com/sunriselayer/sunrise/x/tokenconverter/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	_ "github.com/cosmos/cosmos-sdk/testutil/x/counter" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config"   // import for side-effects
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	icatypes "github.com/cosmos/ibc-go/v9/modules/apps/27-interchain-accounts/types"
	_ "github.com/cosmos/ibc-go/v9/modules/apps/29-fee" // import for side-effects
	ibcfeetypes "github.com/cosmos/ibc-go/v9/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v9/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v9/modules/core/exported"
	_ "github.com/sunriselayer/sunrise/x/shareclass/module"
	shareclassmoduletypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

var (
	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: authtypes.FeeCollectorName},
		{Account: distrtypes.ModuleName},
		{Account: pooltypes.ModuleName},
		{Account: pooltypes.StreamAccount},
		{Account: pooltypes.ProtocolPoolDistrAccount},
		{Account: minttypes.ModuleName, Permissions: []string{authtypes.Minter}},
		{Account: stakingtypes.BondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: stakingtypes.NotBondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: govtypes.ModuleName, Permissions: []string{authtypes.Burner}},
		{Account: nft.ModuleName},
		// IBC module accounts
		{Account: ibctransfertypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: ibcfeetypes.ModuleName},
		{Account: icatypes.ModuleName},
		// this line is used by starport scaffolding # stargate/app/maccPerms

		{Account: damoduletypes.ModuleName},
		{Account: tokenconvertermoduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: liquiditypoolmoduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		{Account: liquidityincentivemoduletypes.ModuleName, Permissions: []string{authtypes.Minter}},
		{Account: swapmoduletypes.ModuleName},
		{Account: feemoduletypes.ModuleName, Permissions: []string{authtypes.Burner}},
		{Account: shareclassmoduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner, authtypes.Staking}},
	}

	// blocked account addresses
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		minttypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		nft.ModuleName,
		// We allow the following module accounts to receive funds:
		// govtypes.ModuleName
		// pooltypes.ModuleName

		damoduletypes.ModuleName,
		tokenconvertermoduletypes.ModuleName,
		liquiditypoolmoduletypes.ModuleName,
		liquidityincentivemoduletypes.ModuleName,
		swapmoduletypes.ModuleName,
		feemoduletypes.ModuleName,
	}

	// application configuration (used by depinject)
	appConfig = appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: runtime.ModuleName,
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName: Name,
					// NOTE: upgrade module is required to be prioritized
					PreBlockers: []string{
						upgradetypes.ModuleName,
						// this line is used by starport scaffolding # stargate/app/preBlockers
					},
					// During begin block slashing happens after distr.BeginBlocker so that
					// there is nothing left over in the validator fee pool, so as to keep the
					// CanWithdrawInvariant invariant.
					// NOTE: staking module is required if HistoricalEntries param > 0
					BeginBlockers: []string{
						minttypes.ModuleName,
						liquidityincentivemoduletypes.ModuleName,
						distrtypes.ModuleName,
						pooltypes.ModuleName,
						slashingtypes.ModuleName,
						evidencetypes.ModuleName,
						stakingtypes.ModuleName,
						authz.ModuleName,
						epochstypes.ModuleName,
						// ibc modules
						ibcexported.ModuleName,
						ibctransfertypes.ModuleName,
						icatypes.ModuleName,
						ibcfeetypes.ModuleName,
						// chain modules
						damoduletypes.ModuleName,
						tokenconvertermoduletypes.ModuleName,
						liquiditypoolmoduletypes.ModuleName,
						swapmoduletypes.ModuleName,
						feemoduletypes.ModuleName,
						selfdelegationmoduletypes.ModuleName,
						shareclassmoduletypes.ModuleName,
						// this line is used by starport scaffolding # stargate/app/beginBlockers
					},
					EndBlockers: []string{
						govtypes.ModuleName,
						stakingtypes.ModuleName,
						feegrant.ModuleName,
						group.ModuleName,
						pooltypes.ModuleName,
						// ibc modules
						ibcexported.ModuleName,
						ibctransfertypes.ModuleName,
						icatypes.ModuleName,
						ibcfeetypes.ModuleName,
						// chain modules
						damoduletypes.ModuleName,
						tokenconvertermoduletypes.ModuleName,
						liquiditypoolmoduletypes.ModuleName,
						liquidityincentivemoduletypes.ModuleName,
						swapmoduletypes.ModuleName,
						feemoduletypes.ModuleName,
						selfdelegationmoduletypes.ModuleName,
						shareclassmoduletypes.ModuleName,
						// this line is used by starport scaffolding # stargate/app/endBlockers
					},
					// The following is mostly only needed when ModuleName != StoreKey name.
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
						{
							ModuleName: accounts.ModuleName,
							KvStoreKey: accounts.StoreKey,
						},
					},
					// NOTE: The genutils module must occur after staking so that pools are
					// properly initialized with tokens from genesis accounts.
					// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
					InitGenesis: []string{
						consensustypes.ModuleName,
						// <sunrise>
						authtypes.ModuleName,
						banktypes.ModuleName,
						stakingtypes.ModuleName,
						tokenconvertermoduletypes.ModuleName,
						selfdelegationmoduletypes.ModuleName,
						shareclassmoduletypes.ModuleName,
						// </sunrise>
						accounts.ModuleName,
						distrtypes.ModuleName,
						slashingtypes.ModuleName,
						govtypes.ModuleName,
						minttypes.ModuleName,
						genutiltypes.ModuleName,
						evidencetypes.ModuleName,
						authz.ModuleName,
						feegrant.ModuleName,
						nft.ModuleName,
						group.ModuleName,
						upgradetypes.ModuleName,
						circuittypes.ModuleName,
						pooltypes.ModuleName,
						epochstypes.ModuleName,
						// ibc modules
						paramstypes.ModuleName,
						ibcexported.ModuleName,
						ibctransfertypes.ModuleName,
						icatypes.ModuleName,
						ibcfeetypes.ModuleName,
						// chain modules
						damoduletypes.ModuleName,
						liquiditypoolmoduletypes.ModuleName,
						liquidityincentivemoduletypes.ModuleName,
						swapmoduletypes.ModuleName,
						feemoduletypes.ModuleName,
						// this line is used by starport scaffolding # stargate/app/initGenesis
					},
					// SkipStoreKeys is an optional list of store keys to skip when constructing the
					// module's keeper. This is useful when a module does not have a store key.
					SkipStoreKeys: []string{
						"tx",
						runtime.ModuleName,
						genutiltypes.ModuleName,
					},
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: appconfig.WrapAny(&authmodulev1.Module{
					Bech32Prefix:             AccountAddressPrefix,
					ModuleAccountPermissions: moduleAccPerms,
					// By default modules authority is the governance module. This is configurable with the following:
					// Authority: "group", // A custom module authority can be set using a module name
					// Authority: "cosmos1cwwv22j5ca08ggdv9c2uky355k908694z577tv", // or a specific address
				}),
			},
			{
				Name: banktypes.ModuleName,
				Config: appconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name:   stakingtypes.ModuleName,
				Config: appconfig.WrapAny(&stakingmodulev1.Module{}),
			},
			{
				Name:   slashingtypes.ModuleName,
				Config: appconfig.WrapAny(&slashingmodulev1.Module{}),
			},
			{
				Name:   "tx",
				Config: appconfig.WrapAny(&txconfigv1.Config{}),
			},
			{
				Name:   genutiltypes.ModuleName,
				Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
			},
			{
				Name:   authz.ModuleName,
				Config: appconfig.WrapAny(&authzmodulev1.Module{}),
			},
			{
				Name:   upgradetypes.ModuleName,
				Config: appconfig.WrapAny(&upgrademodulev1.Module{}),
			},
			{
				Name:   distrtypes.ModuleName,
				Config: appconfig.WrapAny(&distrmodulev1.Module{}),
			},
			{
				Name:   evidencetypes.ModuleName,
				Config: appconfig.WrapAny(&evidencemodulev1.Module{}),
			},
			{
				Name:   minttypes.ModuleName,
				Config: appconfig.WrapAny(&mintmodulev1.Module{}),
			},
			{
				Name: group.ModuleName,
				Config: appconfig.WrapAny(&groupmodulev1.Module{
					MaxExecutionPeriod: durationpb.New(time.Second * 1209600),
					MaxMetadataLen:     255,
				}),
			},
			{
				Name:   nft.ModuleName,
				Config: appconfig.WrapAny(&nftmodulev1.Module{}),
			},
			{
				Name:   feegrant.ModuleName,
				Config: appconfig.WrapAny(&feegrantmodulev1.Module{}),
			},
			{
				Name:   govtypes.ModuleName,
				Config: appconfig.WrapAny(&govmodulev1.Module{}),
			},
			{
				Name:   consensustypes.ModuleName,
				Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
			},
			{
				Name:   circuittypes.ModuleName,
				Config: appconfig.WrapAny(&circuitmodulev1.Module{}),
			},
			{
				Name:   pooltypes.ModuleName,
				Config: appconfig.WrapAny(&poolmodulev1.Module{}),
			},
			{
				Name:   accounts.ModuleName,
				Config: appconfig.WrapAny(&accountsmodulev1.Module{}),
			},
			{
				Name:   epochstypes.ModuleName,
				Config: appconfig.WrapAny(&epochsmodulev1.Module{}),
			},
			{
				Name:   paramstypes.ModuleName,
				Config: appconfig.WrapAny(&paramsmodulev1.Module{}),
			},
			{
				Name:   damoduletypes.ModuleName,
				Config: appconfig.WrapAny(&damoduletypes.Module{}),
			},
			{
				Name:   tokenconvertermoduletypes.ModuleName,
				Config: appconfig.WrapAny(&tokenconvertermoduletypes.Module{}),
			},
			{
				Name:   liquiditypoolmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&liquiditypoolmoduletypes.Module{}),
			},
			{
				Name:   liquidityincentivemoduletypes.ModuleName,
				Config: appconfig.WrapAny(&liquidityincentivemoduletypes.Module{}),
			},
			{
				Name:   swapmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&swapmoduletypes.Module{}),
			},
			{
				Name:   feemoduletypes.ModuleName,
				Config: appconfig.WrapAny(&feemoduletypes.Module{}),
			},
			{
				Name:   selfdelegationmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&selfdelegationmoduletypes.Module{}),
			},
			{
				Name:   shareclassmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&shareclassmoduletypes.Module{}),
			},
			// this line is used by starport scaffolding # stargate/app/moduleConfig
		},
	})
)
