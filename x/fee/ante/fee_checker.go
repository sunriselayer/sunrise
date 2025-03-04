// This file is based on github.com/cosmos/cosmos-sdk/x/auth/ante/validator_tx_fee.go
package ante

import (
	"math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	feekeeper "github.com/sunriselayer/sunrise/x/fee/keeper"
)

func checkTxFeeWithValidatorMinGasPrices(k feekeeper.Keeper) func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	return func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
		return _checkTxFeeWithValidatorMinGasPrices(ctx, tx, k)
	}
}

// checkTxFeeWithValidatorMinGasPrices implements the default fee logic, where the minimum price per
// unit of gas is fixed and set by each validator, can the tx priority is computed from the gas price.
func _checkTxFeeWithValidatorMinGasPrices(ctx sdk.Context, tx sdk.Tx, k feekeeper.Keeper) (sdk.Coins, int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, 0, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	// <sunrise>
	if ctx.BlockHeight() != 0 {
		if len(feeCoins) != 1 {
			return nil, 0, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "only one fee denomination is allowed")
		}
		params := k.GetParams(ctx)
		if feeCoins[0].Denom != params.FeeDenom {
			includedBypass := false
			for _, denom := range params.BypassDenoms {
				if feeCoins[0].Denom == denom {
					includedBypass = true
					break
				}
			}

			if !includedBypass {
				return nil, 0, errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid fee denomination: %s", feeCoins[0].Denom)
			}
		}
	}
	// </sunrise>

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.ExecMode() == sdk.ExecModeCheck { // NOTE: using environment here breaks the API of fee logic, an alternative must be found for server/v2. ref: https://github.com/cosmos/cosmos-sdk/issues/19640
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdkmath.LegacyNewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !feeCoins.IsAnyGTE(requiredFees) {
				return nil, 0, errorsmod.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}

	priority := getTxPriority(feeCoins, int64(gas))
	return feeCoins, priority, nil
}

// getTxPriority returns a naive tx priority based on the amount of the smallest denomination of the gas price
// provided in a transaction.
// NOTE: This implementation should be used with a great consideration as it opens potential attack vectors
// where txs with multiple coins could not be prioritize as expected.
func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}
