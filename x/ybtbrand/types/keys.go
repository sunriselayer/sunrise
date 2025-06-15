package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "ybtbrand"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// State key prefixes
var (
	ParamsKey             = collections.NewPrefix("p_ybtbrand")
	TokensKey             = collections.NewPrefix([]byte{0x02})
	YieldIndexKey         = collections.NewPrefix([]byte{0x03})
	UserLastYieldIndexKey = collections.NewPrefix([]byte{0x04})
)
