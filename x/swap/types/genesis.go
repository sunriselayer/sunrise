package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		IncomingInFlightPacketList: []IncomingInFlightPacket{},
		OutgoingInFlightPacketList: []OutgoingInFlightPacket{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in incomingPacket
	incomingPacketIndexMap := make(map[string]struct{})

	for _, elem := range gs.IncomingInFlightPacketList {
		index := string(IncomingInFlightPacketKey(elem.Index))
		if _, ok := incomingPacketIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for incomingPacket")
		}
		incomingPacketIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in outgoingInFlightPacket
	outgoingInFlightPacketIndexMap := make(map[string]struct{})

	for _, elem := range gs.OutgoingInFlightPacketList {
		index := string(OutgoingInFlightPacketKey(elem.Index))
		if _, ok := outgoingInFlightPacketIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for outgoingInFlightPacket")
		}
		outgoingInFlightPacketIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
