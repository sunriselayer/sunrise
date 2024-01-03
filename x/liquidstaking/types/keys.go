package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "liquidstaking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_liquidstaking"
)

var (
	ParamsKey = []byte("p_liquidstaking")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	// Keys for store prefixes
	LiquidValidatorsKey = []byte{0xc0} // prefix for each key to a liquid validator
)

// GetLiquidValidatorKey creates the key for the liquid validator with address
// VALUE: liquidstaking/LiquidValidator
func GetLiquidValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(LiquidValidatorsKey, address.MustLengthPrefix(operatorAddr)...)
}
