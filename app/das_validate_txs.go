package app

import (
	"encoding/json"

	"cosmossdk.io/log"
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	coretypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/pkg/blob"
)

// separateTxs decodes raw tendermint txs into normal and blob txs.
func separateTxs(_ client.TxConfig, rawTxs [][]byte) ([][]byte, []blob.BlobTx) {
	normalTxs := make([][]byte, 0, len(rawTxs))
	blobTxs := make([]blob.BlobTx, 0, len(rawTxs))
	for _, rawTx := range rawTxs {
		bTx, isBlob := blob.UnmarshalBlobTx(rawTx)
		if isBlob {
			blobTxs = append(blobTxs, bTx)
		} else {
			normalTxs = append(normalTxs, rawTx)
		}
	}
	return normalTxs, blobTxs
}

// FilterTxs applies the antehandler to all proposed transactions and removes transactions that return an error.
func FilterTxs(logger log.Logger, ctx sdk.Context, handler sdk.AnteHandler, txConfig client.TxConfig, txs [][]byte) [][]byte {
	normalTxs, blobTxs := separateTxs(txConfig, txs)
	normalTxs, ctx = filterStdTxsAndVoteExtTx(logger, txConfig.TxDecoder(), ctx, handler, normalTxs)
	blobTxs, _ = filterBlobTxs(logger, txConfig.TxDecoder(), ctx, handler, blobTxs)
	return append(normalTxs, encodeBlobTxs(blobTxs)...)
}

// filterStdTxsAndVoteExtTx collects vote extension tx and applies the provided antehandler to each transaction and removes
// transactions that return an error. Panics are caught by the checkTxValidity
// function used to apply the ante handler.
func filterStdTxsAndVoteExtTx(logger log.Logger, dec sdk.TxDecoder, ctx sdk.Context, handler sdk.AnteHandler, txs [][]byte) ([][]byte, sdk.Context) {
	n := 0
	for _, tx := range txs {
		// collect vote extension tx
		var voteExtTx VoteExtensionTx
		if err := json.Unmarshal(tx, &voteExtTx); err == nil {
			txs[n] = tx
			n++
			continue
		}

		sdkTx, err := dec(tx)
		if err != nil {
			logger.Error("decoding already checked transaction", "tx", tmbytes.HexBytes(coretypes.Tx(tx).Hash()), "error", err)
			continue
		}
		ctx, err = handler(ctx, sdkTx, false)
		// either the transaction is invalid (ie incorrect nonce) and we
		// simply want to remove this tx, or we're catching a panic from one
		// of the anteHanders which is logged.
		if err != nil {
			logger.Error(
				"filtering already checked transaction",
				"tx", tmbytes.HexBytes(coretypes.Tx(tx).Hash()),
				"error", err,
				"msgs", msgTypes(sdkTx),
			)
			telemetry.IncrCounter(1, "prepare_proposal", "invalid_std_txs")
			continue
		}
		txs[n] = tx
		n++

	}

	return txs[:n], ctx
}

// filterBlobTxs applies the provided antehandler to each transaction
// and removes transactions that return an error. Panics are caught by the checkTxValidity
// function used to apply the ante handler.
func filterBlobTxs(logger log.Logger, dec sdk.TxDecoder, ctx sdk.Context, handler sdk.AnteHandler, txs []blob.BlobTx) ([]blob.BlobTx, sdk.Context) {
	n := 0
	for _, tx := range txs {
		sdkTx, err := dec(tx.Tx)
		if err != nil {
			logger.Error("decoding already checked blob transaction", "tx", tmbytes.HexBytes(coretypes.Tx(tx.Tx).Hash()), "error", err)
			continue
		}
		ctx, err = handler(ctx, sdkTx, false)
		// either the transaction is invalid (ie incorrect nonce) and we
		// simply want to remove this tx, or we're catching a panic from one
		// of the anteHanders which is logged.
		if err != nil {
			logger.Error(
				"filtering already checked blob transaction", "tx", tmbytes.HexBytes(coretypes.Tx(tx.Tx).Hash()), "error", err,
			)
			telemetry.IncrCounter(1, "prepare_proposal", "invalid_blob_txs")
			continue
		}
		txs[n] = tx
		n++

	}

	return txs[:n], ctx
}

func msgTypes(sdkTx sdk.Tx) []string {
	msgs := sdkTx.GetMsgs()
	msgNames := make([]string, len(msgs))
	for i, msg := range msgs {
		msgNames[i] = sdk.MsgTypeURL(msg)
	}
	return msgNames
}

func encodeBlobTxs(blobTxs []blob.BlobTx) [][]byte {
	txs := make([][]byte, len(blobTxs))
	var err error
	for i, tx := range blobTxs {
		txs[i], err = blob.MarshalBlobTx(tx.Tx, tx.Blobs...)
		if err != nil {
			panic(err)
		}
	}
	return txs
}
