package app

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/sunrise-zone/sunrise-app/pkg/blob"
)

// FinalizeBlock will execute the block proposal provided by RequestFinalizeBlock.
// Specifically, it will execute an application's BeginBlock (if defined), followed
// by the transactions in the proposal, finally followed by the application's
// EndBlock (if defined).
//
// For each raw transaction, i.e. a byte slice, BaseApp will only execute it if
// it adheres to the sdk.Tx interface. Otherwise, the raw transaction will be
// skipped. This is to support compatibility with proposers injecting vote
// extensions into the proposal, which should not themselves be executed in cases
// where they adhere to the sdk.Tx interface.
func (app *App) FinalizeBlock(req *abci.RequestFinalizeBlock) (*abci.ResponseFinalizeBlock, error) {
	txs := [][]byte{}
	// Iterate over all raw transactions in the proposal and attempt to execute
	// them, gathering the execution results.
	//
	// NOTE: Not all raw transactions may adhere to the sdk.Tx interface, e.g.
	// vote extensions, so skip those.
	for _, rawTx := range req.Txs {
		tx := rawTx
		blobTx, isBlobTx := blob.UnmarshalBlobTx(rawTx)
		if isBlobTx {
			tx = blobTx.Tx
		}
		txs = append(txs, tx)
	}
	req.Txs = txs

	return app.BaseApp.FinalizeBlock(req)
}
