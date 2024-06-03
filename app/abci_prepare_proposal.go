package app

import (
	"time"

	"github.com/sunriselayer/sunrise/app/ante"
	"github.com/sunriselayer/sunrise/pkg/da"
	"github.com/sunriselayer/sunrise/pkg/shares"
	"github.com/sunriselayer/sunrise/pkg/square"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/telemetry"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewProposalContext returns a context with a branched version of the state
// that is safe to query during ProcessProposal.
func (app *App) NewProposalContext(header cmtproto.Header) sdk.Context {
	// use custom query multistore if provided
	ms := app.CommitMultiStore().CacheMultiStore()
	ctx := sdk.NewContext(ms, header, false, app.Logger()).WithBlockGasMeter(storetypes.NewInfiniteGasMeter())
	ctx = ctx.WithConsensusParams(app.GetConsensusParams(ctx))

	return ctx
}

// PrepareProposal fulfills the sunrise-core version of the ABCI interface by
// preparing the proposal block data. The square size is determined by first
// estimating it via the size of the passed block data. Then, this method
// generates the data root for the proposal block and passes it back to
// tendermint via the BlockData. Panics indicate a developer error and should
// immediately halt the node for visibility and so they can be quickly resolved.
func (app *App) PrepareProposal(req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
	defer telemetry.MeasureSince(time.Now(), "prepare_proposal")

	// create a context using a branch of the state and loaded using the
	// proposal height and chain-id
	sdkCtx := app.NewProposalContext(cmtproto.Header{
		ChainID: app.ChainID(),
		Height:  req.Height,
		Time:    req.Time,
	})
	// filter out invalid transactions.
	// TODO: we can remove all state independent checks from the ante handler here such as signature verification
	// and only check the state dependent checks like fees and nonces as all these transactions have already
	// passed CheckTx.
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

	var txs [][]byte
	// This if statement verifies whether the preparation of the proposal
	// pertains to the first block. If it does, the block is constructed using
	// an empty set of transactions. However, even without this validation,
	// the initial block is anticipated to be devoid of transactions, as
	// established by the findings presented in
	// https://github.com/sunriselayer/sunrise/issues/1899;
	// The inclusion of this check is out of an abundance of caution.
	// The rationale behind having an empty first block revolves around the fact
	// that no transactions can enter the mempool since no committed state exists
	// until after the first block is committed (at which point the Genesis state
	// gets committed too). Consequently, the prepare proposal request for the
	// first block is expected to contain no transaction, so is the first block.
	if app.LastBlockHeight() == 0 {
		txs = make([][]byte, 0)
	} else {
		txs = FilterTxs(app.Logger(), sdkCtx, handler, app.txConfig, req.Txs)
	}

	// build the square from the set of valid and prioritised transactions.
	// The txs returned are the ones used in the square and block
	dataSquare, txs, err := square.Build(txs, app.BaseApp.AppVersion(), app.GovSquareSizeUpperBound(sdkCtx))
	if err != nil {
		panic(err)
	}

	// erasure the data square which we use to create the data root.
	// Note: uses the nmt wrapper to construct the tree.
	// checkout pkg/wrapper/nmt_wrapper.go for more information.
	eds, err := da.ExtendShares(shares.ToBytes(dataSquare))
	if err != nil {
		app.Logger().Error(
			"failure to erasure the data square while creating a proposal block",
			"error",
			err.Error(),
		)
		panic(err)
	}

	// create the new data root by creating the data availability header (merkle
	// roots of each row and col of the erasure data).
	dah, err := da.NewDataAvailabilityHeader(eds)
	if err != nil {
		app.Logger().Error(
			"failure to create new data availability header",
			"error",
			err.Error(),
		)
		panic(err)
	}

	// tendermint doesn't need to use any of the erasure data, as only the
	// protobuf encoded version of the block data is gossiped.
	return &abci.ResponsePrepareProposal{
		Txs:        txs,
		DataHash:   dah.Hash(),
		SquareSize: uint64(dataSquare.Size()),
	}, nil
}
