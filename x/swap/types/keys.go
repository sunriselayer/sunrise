package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "swap"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("params")

var (
	IncomingInFlightPacketsKey = collections.NewPrefix("incoming_in_flight_packets")
	OutgoingInFlightPacketsKey = collections.NewPrefix("outgoing_in_flight_packets")
)

// deprecated
func KeyPrefix(p string) []byte {
	return []byte(p)
}
