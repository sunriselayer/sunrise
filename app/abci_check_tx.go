package app

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/skip-mev/block-sdk/v2/abci/checktx"
	"github.com/sunriselayer/sunrise/pkg/blob"
	blobtypes "github.com/sunriselayer/sunrise/x/blob/types"
)

// CheckTx implements the ABCI interface and executes a tx in CheckTx mode. This
// method wraps the default Baseapp's method so that it can parse and check
// transactions that contain blobs.
func (app *App) CheckTx(req *abci.RequestCheckTx) (*abci.ResponseCheckTx, error) {
	tx := req.Tx
	// check if the transaction contains blobs
	btx, isBlob := blob.UnmarshalBlobTx(tx)

	if !isBlob {
		// reject transactions that can't be decoded
		sdkTx, err := app.txConfig.TxDecoder()(tx)
		if err != nil {
			return sdkerrors.ResponseCheckTxWithEvents(err, 0, 0, []abci.Event{}, false), err
		}
		// reject transactions that have a MsgPFB but no blobs attached to the tx
		for _, msg := range sdkTx.GetMsgs() {
			if _, ok := msg.(*blobtypes.MsgPayForBlobs); !ok {
				continue
			}
			return sdkerrors.ResponseCheckTxWithEvents(blobtypes.ErrNoBlobs, 0, 0, []abci.Event{}, false), blobtypes.ErrNoBlobs
		}
		// don't do anything special if we have a normal transaction
		return app.CheckTxHandler(req)
	}

	switch req.Type {
	// new transactions must be checked in their entirety
	case abci.CheckTxType_New:
		err := blobtypes.ValidateBlobTx(app.txConfig, btx)
		if err != nil {
			return sdkerrors.ResponseCheckTxWithEvents(err, 0, 0, []abci.Event{}, false), err
		}
	case abci.CheckTxType_Recheck:
	default:
		panic(fmt.Sprintf("unknown RequestCheckTx type: %s", req.Type))
	}

	req.Tx = btx.Tx
	return app.CheckTxHandler(req)
}

// SetCheckTx sets the CheckTxHandler for the app.
func (app *App) SetCheckTx(handler checktx.CheckTx) {
	app.CheckTxHandler = handler
}
