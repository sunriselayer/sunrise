package v0_3_0_test

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/sunriselayer/sunrise/app/upgrades"
	datypes "github.com/sunriselayer/sunrise/x/da/types"
)

const UpgradeName string = "v0_3_0_test"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{datypes.StoreKey},
		Deleted: []string{datypes.StoreKey},
	},
}
