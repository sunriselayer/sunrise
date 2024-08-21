package types

import (
	"fmt"

	// "github.com/sunriselayer/sunrise/pkg/appconsts"

	"cosmossdk.io/errors"
)

// DefaultParamspace defines the default blobstream module parameter subspace
const (
	DefaultParamspace = ModuleName

	// MinimumDataCommitmentWindow is a constant that defines the minimum
	// allowable window for the Blobstream data commitments.
	MinimumDataCommitmentWindow = 100
)

// ParamsStoreKeyDataCommitmentWindow is the key used for the
// DataCommitmentWindow param.
var ParamsStoreKeyDataCommitmentWindow = []byte("DataCommitmentWindow")

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: Params{
			DataCommitmentWindow: 400,
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	if err := gs.Params.Validate(); err != nil {
		return errors.Wrap(err, "params")
	}
	return nil
}

func validateDataCommitmentWindow(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < MinimumDataCommitmentWindow {
		return errors.Wrap(ErrInvalidDataCommitmentWindow, fmt.Sprintf(
			"data commitment window %v must be >= minimum data commitment window %v",
			val,
			MinimumDataCommitmentWindow,
		))
	}
	// if val > uint64(appconsts.DataCommitmentBlocksLimit) {
	// 	return errors.Wrap(ErrInvalidDataCommitmentWindow, fmt.Sprintf(
	// 		"data commitment window %v must be <= data commitment blocks limit %v",
	// 		val,
	// 		appconsts.DataCommitmentBlocksLimit,
	// 	))
	// }
	return nil
}
