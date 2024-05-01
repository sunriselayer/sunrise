package types

import (
	fmt "fmt"

	"github.com/sunriselayer/sunrise-app/pkg/appconsts"
	"github.com/sunriselayer/sunrise-app/pkg/shares"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyGasPerBlobByte              = []byte("GasPerBlobByte")
	DefaultGasPerBlobByte   uint32 = appconsts.DefaultGasPerBlobByte
	KeyGovMaxSquareSize            = []byte("GovMaxSquareSize")
	DefaultGovMaxSquareSize uint64 = appconsts.DefaultGovMaxSquareSize
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(gasPerBlobByte uint32, govMaxSquareSize uint64) Params {
	return Params{
		GasPerBlobByte:   gasPerBlobByte,
		GovMaxSquareSize: govMaxSquareSize,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(DefaultGasPerBlobByte, appconsts.DefaultGovMaxSquareSize)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyGasPerBlobByte, &p.GasPerBlobByte, validateGasPerBlobByte),
		paramtypes.NewParamSetPair(KeyGovMaxSquareSize, &p.GovMaxSquareSize, validateGovMaxSquareSize),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	err := validateGasPerBlobByte(p.GasPerBlobByte)
	if err != nil {
		return err
	}
	return validateGovMaxSquareSize(p.GovMaxSquareSize)
}

// validateGasPerBlobByte validates the GasPerBlobByte param
func validateGasPerBlobByte(v interface{}) error {
	gasPerBlobByte, ok := v.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if gasPerBlobByte == 0 {
		return fmt.Errorf("gas per blob byte cannot be 0")
	}

	return nil
}

// validateGovMaxSquareSize validates the GovMaxSquareSize param
func validateGovMaxSquareSize(v interface{}) error {
	govMaxSquareSize, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if govMaxSquareSize == 0 {
		return fmt.Errorf("gov max square size cannot be zero")
	}

	if !shares.IsPowerOfTwo(govMaxSquareSize) {
		return fmt.Errorf(
			"gov max square size must be a power of two: %d",
			govMaxSquareSize,
		)
	}

	return nil
}
