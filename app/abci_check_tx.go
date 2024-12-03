package app

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/skip-mev/block-sdk/v2/abci/checktx"
)

// CheckTx implements the ABCI interface and executes a tx in CheckTx mode. This
// method wraps the default Baseapp's method so that it can parse and check
// transactions that contain blobs.
func (app *App) CheckTx(req *abci.RequestCheckTx) (*abci.ResponseCheckTx, error) {
	return app.CheckTxHandler(req)
}

// SetCheckTx sets the CheckTxHandler for the app.
func (app *App) SetCheckTx(handler checktx.CheckTx) {
	app.CheckTxHandler = handler
}
