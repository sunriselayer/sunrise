package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		InFlightPacketList: []InFlightPacket{},
		AckWaitingPacketList: []AckWaitingPacket{},
// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in inFlightPacket
	inFlightPacketIndexMap := make(map[string]struct{})

	for _, elem := range gs.InFlightPacketList {
		index := string(InFlightPacketKey(elem.Index))
		if _, ok := inFlightPacketIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for inFlightPacket")
		}
		inFlightPacketIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in ackWaitingPacket
ackWaitingPacketIndexMap := make(map[string]struct{})

for _, elem := range gs.AckWaitingPacketList {
	index := string(AckWaitingPacketKey(elem.Index))
	if _, ok := ackWaitingPacketIndexMap[index]; ok {
		return fmt.Errorf("duplicated index for ackWaitingPacket")
	}
	ackWaitingPacketIndexMap[index] = struct{}{}
}
// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
