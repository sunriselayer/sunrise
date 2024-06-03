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
	// Check for duplicated index in ackWaitingPacket
	ackWaitingPacketIndexMap := make(map[string]struct{})

	for _, elem := range gs.IncomingInFlightPacketList {
		index := string(IncomingInFlightPacketKey(elem.Index))
		if _, ok := ackWaitingPacketIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for ackWaitingPacket")
		}
		ackWaitingPacketIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in inFlightPacket
	inFlightPacketIndexMap := make(map[string]struct{})

	for _, elem := range gs.OutgoingInFlightPacketList {
		index := string(OutgoingInFlightPacketKey(elem.Index))
		if _, ok := inFlightPacketIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for inFlightPacket")
		}
		inFlightPacketIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
