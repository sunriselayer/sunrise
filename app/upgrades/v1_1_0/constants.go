// This file defines the upgrade name for the v1.1.0 upgrade.
package v1_1_0

import (
	storetypes "cosmossdk.io/store/types"
	ibchookstypes "github.com/sunriselayer/ibc-hooks/v10/types"
)

// UpgradeName is the name of the upgrade.
const UpgradeName = "v1.1.0"

// StoreUpgrades defines the store upgrades for the v1.1.0 upgrade.
var StoreUpgrades = storetypes.StoreUpgrades{
	Added: []string{
		ibchookstypes.StoreKey,
	},
	Deleted: []string{},
}
