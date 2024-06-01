package app

import (
	"bytes"
	"fmt"
	"time"

	"github.com/sunriselayer/sunrise/app/ante"
	"github.com/sunriselayer/sunrise/pkg/blob"
	"github.com/sunriselayer/sunrise/pkg/da"
	"github.com/sunriselayer/sunrise/pkg/shares"
	"github.com/sunriselayer/sunrise/pkg/square"

	blobtypes "github.com/sunriselayer/sunrise/x/blob/types"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

const rejectedPropBlockLog = "Rejected proposal block:"

func (app *App) ProcessProposal(req *abci.RequestProcessProposal) (retResp *abci.ResponseProcessProposal, retErr error) {
	defer telemetry.MeasureSince(time.Now(), "process_proposal")
	// In the case of a panic from an unexpected condition, it is better for the liveness of the
	// network that we catch it, log an error and vote nil than to crash the node.
	defer func() {
		if err := recover(); err != nil {
			logInvalidPropBlock(app.Logger(), req.ProposerAddress, fmt.Sprintf("caught panic: %v", err))
			telemetry.IncrCounter(1, "process_proposal", "panics")
			resp, err := reject(fmt.Errorf("%s", err))

			retResp = resp
			retErr = err
		}
	}()

	// Create the anteHander that are used to check the validity of
	// transactions. All transactions need to be equally validated here
	// so that the nonce number is always correctly incremented (which
	// may affect the validity of future transactions).
	handler := ante.NewAnteHandler(
		app.AccountKeeper,
		app.BankKeeper,
		app.FeeGrantKeeper,
		app.BlobKeeper,
		app.FeeKeeper,
		app.txConfig.SignModeHandler(),
		ante.DefaultSigVerificationGasConsumer,
		app.IBCKeeper,
	)
	sdkCtx := app.NewProposalContext(cmtproto.Header{
		ChainID: app.ChainID(),
		Height:  req.Height,
		Time:    req.Time,
	})

	if len(req.Txs) == 0 {
		return accept()
	}

	txs := req.Txs

	// iterate over all txs and ensure that all blobTxs are valid, PFBs are correctly signed and non
	// blobTxs have no PFBs present
	for idx, rawTx := range txs {
		tx := rawTx
		blobTx, isBlobTx := blob.UnmarshalBlobTx(rawTx)
		if isBlobTx {
			tx = blobTx.Tx
		}

		sdkTx, err := app.txConfig.TxDecoder()(tx)
		if err != nil {
			// we don't reject the block here because it is not a block validity
			// rule that all transactions included in the block data are
			// decodable
			continue
		}

		// handle non-blob transactions first
		if !isBlobTx {
			msgs := sdkTx.GetMsgs()

			_, has := hasPFB(msgs)
			if has {
				// A non blob tx has a PFB, which is invalid
				err := fmt.Errorf("tx %d has PFB but is not a blob tx", idx)
				logInvalidPropBlock(app.Logger(), req.ProposerAddress, err.Error())
				return reject(err)
			}

			// we need to increment the sequence for every transaction so that
			// the signature check below is accurate. this error only gets hit
			// if the account in question doesn't exist.
			sdkCtx, err = handler(sdkCtx, sdkTx, false)
			if err != nil {
				logInvalidPropBlockError(app.Logger(), req.ProposerAddress, "failure to increment sequence", err)
				return reject(err)
			}

			// we do not need to perform further checks on this transaction,
			// since it has no PFB
			continue
		}

		// validate the blobTx. This is the same validation used in CheckTx ensuring
		// - there is one PFB
		// - that each blob has a valid namespace
		// - that the sizes match
		// - that the namespaces match between blob and PFB
		// - that the share commitment is correct
		if err := blobtypes.ValidateBlobTx(app.txConfig, blobTx); err != nil {
			logInvalidPropBlockError(app.Logger(), req.ProposerAddress, fmt.Sprintf("invalid blob tx %d", idx), err)
			return reject(err)
		}

		// validated the PFB signature
		sdkCtx, err = handler(sdkCtx, sdkTx, false)
		if err != nil {
			logInvalidPropBlockError(app.Logger(), req.ProposerAddress, "invalid PFB signature", err)
			return reject(err)
		}
	}

	// Construct the data square from the block's transactions
	dataSquare, err := square.Construct(txs, app.BaseApp.AppVersion(), app.GovSquareSizeUpperBound(sdkCtx))
	if err != nil {
		logInvalidPropBlockError(app.Logger(), req.ProposerAddress, "failure to compute data square from transactions:", err)
		return reject(err)
	}

	// Assert that the square size stated by the proposer is correct
	if uint64(dataSquare.Size()) != req.SquareSize {
		err := fmt.Errorf("proposed square size differs from calculated square size, expected %d, got %d", req.SquareSize, dataSquare.Size())
		logInvalidPropBlock(app.Logger(), req.ProposerAddress, err.Error())
		return reject(err)
	}

	eds, err := da.ExtendShares(shares.ToBytes(dataSquare))
	if err != nil {
		logInvalidPropBlockError(app.Logger(), req.ProposerAddress, "failure to erasure the data square", err)
		return reject(err)
	}

	dah, err := da.NewDataAvailabilityHeader(eds)
	if err != nil {
		logInvalidPropBlockError(app.Logger(), req.ProposerAddress, "failure to create new data availability header", err)
		return reject(err)
	}
	// by comparing the hashes we know the computed IndexWrappers (with the share indexes of the PFB's blobs)
	// are identical and that square layout is consistent. This also means that the share commitment rules
	// have been followed and thus each blobs share commitment should be valid
	if !bytes.Equal(dah.Hash(), req.DataHash) {
		err := fmt.Errorf("proposed data root %X differs from calculated data root %X", req.DataHash, dah.Hash())
		logInvalidPropBlock(app.Logger(), req.ProposerAddress, err.Error())
		return reject(err)
	}

	return accept()
}

func hasPFB(msgs []sdk.Msg) (*blobtypes.MsgPayForBlobs, bool) {
	for _, msg := range msgs {
		if pfb, ok := msg.(*blobtypes.MsgPayForBlobs); ok {
			return pfb, true
		}
	}
	return nil, false
}

func logInvalidPropBlock(l log.Logger, proposerAddress []byte, reason string) {
	l.Error(
		rejectedPropBlockLog,
		"reason",
		reason,
		"proposer",
		proposerAddress,
	)
}

func logInvalidPropBlockError(l log.Logger, proposerAddress []byte, reason string, err error) {
	l.Error(
		rejectedPropBlockLog,
		"reason",
		reason,
		"proposer",
		proposerAddress,
		"err",
		err.Error(),
	)
}

func reject(err error) (*abci.ResponseProcessProposal, error) {
	return &abci.ResponseProcessProposal{
		Status: abci.ResponseProcessProposal_REJECT,
	}, err
}

func accept() (*abci.ResponseProcessProposal, error) {
	return &abci.ResponseProcessProposal{
		Status: abci.ResponseProcessProposal_ACCEPT,
	}, nil
}
