// This file defines the upgrade name for the v1.1.0 upgrade.
package v1_2_0

import (
	storetypes "cosmossdk.io/store/types"
	crontypes "github.com/sunriselayer/sunrise/x/cron/types"
	tokenfactorytypes "github.com/sunriselayer/sunrise/x/tokenfactory/types"
)

// UpgradeName is the name of the upgrade.
const UpgradeName = "v1.2.0"

// StoreUpgrades defines the store upgrades for the v1.1.0 upgrade.
var StoreUpgrades = storetypes.StoreUpgrades{
	Added: []string{
		tokenfactorytypes.StoreKey,
		crontypes.StoreKey,
	},
	Deleted: []string{},
}
