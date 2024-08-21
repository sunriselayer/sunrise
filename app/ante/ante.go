package ante

import (
	// blobante "github.com/sunriselayer/sunrise/x/blob/ante"
	blob "github.com/sunriselayer/sunrise/x/blob/keeper"
	fee "github.com/sunriselayer/sunrise/x/fee/keeper"

	"cosmossdk.io/x/tx/signing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	ibcante "github.com/cosmos/ibc-go/v8/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	auctionante "github.com/skip-mev/block-sdk/v2/x/auction/ante"
	auctionkeeper "github.com/skip-mev/block-sdk/v2/x/auction/keeper"
	feeante "github.com/sunriselayer/sunrise/x/fee/ante"
)

func NewAnteHandler(
	accountKeeper ante.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	feegrantKeeper ante.FeegrantKeeper,
	blobKeeper blob.Keeper,
	feeKeeper fee.Keeper,
	signModeHandler *signing.HandlerMap,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	channelKeeper *ibckeeper.Keeper,
	auctionkeeper auctionkeeper.Keeper,
	MEVLane auctionante.MEVLane,
	TxEncoder sdk.TxEncoder,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		// Wraps the panic with the string format of the transaction
		NewHandlePanicDecorator(),
		// Set up the context with a gas meter.
		// Must be called before gas consumption occurs in any other decorator.
		ante.NewSetUpContextDecorator(),
		// Ensure the tx does not contain any extension options.
		ante.NewExtensionOptionsDecorator(nil),
		// Ensure the tx passes ValidateBasic.
		ante.NewValidateBasicDecorator(),
		// Ensure the tx has not reached a height timeout.
		ante.NewTxTimeoutHeightDecorator(),
		// Ensure the tx memo <= max memo characters.
		ante.NewValidateMemoDecorator(accountKeeper),
		// Ensure the tx's gas limit is > the gas consumed based on the tx size.
		// Side effect: consumes gas from the gas meter.
		ante.NewConsumeGasForTxSizeDecorator(accountKeeper),
		// Ensure the feepayer (fee granter or first signer) has enough funds to pay for the tx.
		// Side effect: deducts fees from the fee payer. Sets the tx priority in context.
		// ante.NewDeductFeeDecorator(accountKeeper, bankKeeper, feegrantKeeper, nil),
		feeante.NewDeductFeeDecorator(accountKeeper, bankKeeper, feegrantKeeper, feeKeeper),
		// Set public keys in the context for fee-payer and all signers.
		// Contract: must be called before all signature verification decorators.
		ante.NewSetPubKeyDecorator(accountKeeper),
		// Ensure that the tx's count of signatures is <= the tx signature limit.
		ante.NewValidateSigCountDecorator(accountKeeper),
		// Ensure that the tx's gas limit is > the gas consumed based on signature verification.
		// Side effect: consumes gas from the gas meter.
		ante.NewSigGasConsumeDecorator(accountKeeper, sigGasConsumer),
		// Ensure that the tx's signatures are valid. For each signature, ensure
		// that the signature's sequence number (a.k.a nonce) matches the
		// account sequence number of the signer.
		// Note: does not consume gas from the gas meter.
		ante.NewSigVerificationDecorator(accountKeeper, signModeHandler),
		// Ensure that the tx's gas limit is > the gas consumed based on the blob size(s).
		// Contract: must be called after all decorators that consume gas.
		// Note: does not consume gas from the gas meter.
		// blobante.NewMinGasPFBDecorator(blobKeeper),
		// Ensure that the tx's total blob size is <= the max blob size.
		// blobante.NewMaxBlobSizeDecorator(blobKeeper),
		// Side effect: increment the nonce for all tx signers.
		ante.NewIncrementSequenceDecorator(accountKeeper),
		// Ensure that the tx is not a IBC packet or update message that has already been processed.
		ibcante.NewRedundantRelayDecorator(channelKeeper),
		auctionante.NewAuctionDecorator(auctionkeeper, TxEncoder, MEVLane),
	)
}

var DefaultSigVerificationGasConsumer = ante.DefaultSigVerificationGasConsumer
