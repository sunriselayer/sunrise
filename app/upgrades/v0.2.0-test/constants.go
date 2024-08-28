package v0_2_0_test

import (
	storetypes "cosmossdk.io/store/types"
	blobtypes "github.com/sunriselayer/sunrise/x/blob/types"
	bstypes "github.com/sunriselayer/sunrise/x/blobstream/types"
	datypes "github.com/sunriselayer/sunrise/x/da/types"

	"github.com/sunriselayer/sunrise/app/upgrades"
)

const UpgradeName string = "v0_2_0_test"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{datypes.ModuleName},
		Deleted: []string{blobtypes.ModuleName, bstypes.ModuleName},
	},
}
