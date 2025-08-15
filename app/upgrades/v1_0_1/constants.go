// This file defines the upgrade name for the v1.0.1 upgrade.
package v1_0_1

import (
	storetypes "cosmossdk.io/store/types"
	ibchookstypes "github.com/sunriselayer/ibc-hooks/v10/types"
)

// UpgradeName is the name of the upgrade.
const UpgradeName = "v1.0.1"

// StoreUpgrades defines the store upgrades for the v1.0.1 upgrade.
var StoreUpgrades = storetypes.StoreUpgrades{
	Added: []string{
		ibchookstypes.StoreKey,
	},
	Deleted: []string{},
}
