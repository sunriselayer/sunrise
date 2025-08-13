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
