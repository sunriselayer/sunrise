# Upgrades

This directory manages the application's upgrade handlers. Each upgrade should have its own dedicated directory, named after the version.

## How to Add a New Upgrade

Follow these steps to add a new upgrade:

1.  **Create a New Directory**:
    Inside `app/upgrades`, create a new directory for the new version. For example, for an upgrade to `v1.1.0`, create the directory `app/upgrades/v1_1_0`.

2.  **Define the Upgrade Name**:
    In the newly created directory, create a `constants.go` file and define the upgrade name as a constant.

    ```go
    // app/upgrades/v1_1_0/constants.go
    package v1_1_0

    // UpgradeName is the name of the upgrade.
    const UpgradeName = "v1.1.0"
    ```

3.  **Create the Upgrade Handler**:
    In the same directory, create an `upgrade.go` file. This file will implement the `UpgradeHandler` logic.

    ```go
    // app/upgrades/v1_1_0/upgrade.go
    package v1_1_0

    import (
    	"context"
    	"fmt"

    	"cosmossdk.io/x/upgrade/types"
    	sdk "github.com/cosmos/cosmos-sdk/types"
    	"github.com/cosmos/cosmos-sdk/types/module"
    )

    // CreateUpgradeHandler creates a handler for the v1.1.0 upgrade.
    func CreateUpgradeHandler(
    	mm *module.Manager,
    	configurator module.Configurator,
    ) types.UpgradeHandler {
    	return func(goCtx context.Context, plan types.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
    		ctx := sdk.UnwrapSDKContext(goCtx)
    		ctx.Logger().Info(fmt.Sprintf("starting upgrade to %s...", UpgradeName))

    		// Run migrations.
    		return mm.RunMigrations(goCtx, configurator, fromVM)
    	}
    }
    ```

4.  **Register the Handler in `app.go`**:
    Finally, register the new upgrade handler in the `setupUpgradeHandlers` function in `app.go`.

    First, import the new upgrade package.

    ```go
    // app.go
    import (
        // ...
        v1_0_1 "github.com/sunriselayer/sunrise/app/upgrades/v1_0_1"
        v1_1_0 "github.com/sunriselayer/sunrise/app/upgrades/v1_1_0" // Add this line
    )
    ```

    Next, set up the handler in `setupUpgradeHandlers`.

    ```go
    // app.go
    func (app *App) setupUpgradeHandlers() {
        // ... existing handlers
        app.UpgradeKeeper.SetUpgradeHandler(
            v1_1_0.UpgradeName,
            v1_1_0.CreateUpgradeHandler(app.ModuleManager, app.Configurator()),
        )
    }
    ```

5.  **Store Migrations (if necessary)**:
    If the upgrade requires store changes, such as adding or deleting a module, you need to define `StoreUpgrades`. This should be done in `constants.go` alongside the `UpgradeName`.

    In `app/upgrades/v1_1_0/constants.go`, define the `StoreUpgrades` variable:

    ```go
    // app/upgrades/v1_1_0/constants.go
    package v1_1_0

    import (
    	storetypes "cosmossdk.io/store/types"
    	// newmoduletypes "path/to/new/module/types"
    )

    // UpgradeName is the name of the upgrade.
    const UpgradeName = "v1.1.0"

    // StoreUpgrades defines the store upgrades for this version.
    var StoreUpgrades = storetypes.StoreUpgrades{
    	Added: []string{
    		// newmoduletypes.StoreKey,
    	},
    	Deleted: []string{
    		// "oldmodule",
    	},
    }
    ```

    In `app.go`, modify the `setupUpgradeHandlers` function to use the `StoreUpgrades` variable. When adding multiple upgrades, this logic can be extended with a `switch` statement.

    ```go
    // app.go
    func (app *App) setupUpgradeHandlers() {
        // ... existing handlers are registered here, for instance
        app.UpgradeKeeper.SetUpgradeHandler(
            v1_1_0.UpgradeName,
            v1_1_0.CreateUpgradeHandler(app.ModuleManager, app.Configurator()),
        )
        // ...

        upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
        if err != nil {
            panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
        }

        if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
            return
        }

        if upgradeInfo.Name == v1_1_0.UpgradeName {
            // Note: in app.go, upgradetypes is imported as "cosmossdk.io/x/upgrade/types"
            app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &v1_1_0.StoreUpgrades))
        }
    }
    ```
