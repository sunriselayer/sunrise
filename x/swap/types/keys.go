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

var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("params/")

	IncomingInFlightPacketsKeyPrefix = collections.NewPrefix("incoming_in_flight_packets/")
	OutgoingInFlightPacketsKeyPrefix = collections.NewPrefix("outgoing_in_flight_packets/")
)

var (
	IncomingInFlightPacketsKeyCodec = collections.TripleKeyCodec(collections.StringKey, collections.StringKey, collections.Uint64Key)
	OutgoingInFlightPacketsKeyCodec = collections.TripleKeyCodec(collections.StringKey, collections.StringKey, collections.Uint64Key)
)

func IncomingInFlightPacketKey(index PacketIndex) collections.Triple[string, string, uint64] {
	return collections.Join3(index.PortId, index.ChannelId, index.Sequence)
}

func OutgoingInFlightPacketKey(index PacketIndex) collections.Triple[string, string, uint64] {
	return collections.Join3(index.PortId, index.ChannelId, index.Sequence)
}
