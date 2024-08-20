package v0_1_5_test

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/sunriselayer/sunrise/app/upgrades"
)

const UpgradeName string = "v0_1_5_test"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
