package v0_2_3_test

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/sunriselayer/sunrise/app/upgrades"
)

const UpgradeName string = "v0_2_3_test"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
