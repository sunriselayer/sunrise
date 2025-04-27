// https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-rc.1/simapp/ante.go
package app

import (
	errorsmod "cosmossdk.io/errors"

	ibcante "github.com/cosmos/ibc-go/v10/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante/unorderedtx"

	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	circuitante "cosmossdk.io/x/circuit/ante"
	feeante "github.com/sunriselayer/sunrise/x/fee/ante"
	fee "github.com/sunriselayer/sunrise/x/fee/keeper"
)

// HandlerOptions are the options required for constructing a default SDK AnteHandler.
type HandlerOptions struct {
	ante.HandlerOptions

	IBCKeeper     *keeper.Keeper
	CircuitKeeper *circuitante.CircuitBreaker
	FeeKeeper     *fee.Keeper
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	if options.IBCKeeper == nil {
		return nil, errors.New("ibc keeper is required for ante builder")
	}

	if options.CircuitKeeper == nil {
		return nil, errors.New("circuit keeper is required for ante builder")
	}

	if options.FeeKeeper == nil {
		return nil, errors.New("fee keeper is required for ante builder")
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		feeante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.FeeKeeper),
		// NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
