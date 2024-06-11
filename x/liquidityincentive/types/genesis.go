package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		EpochList: []Epoch{},
		GaugeList: []Gauge{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in epoch
	epochIdMap := make(map[uint64]bool)
	epochCount := gs.GetEpochCount()
	for _, elem := range gs.EpochList {
		if _, ok := epochIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for epoch")
		}
		if elem.Id >= epochCount {
			return fmt.Errorf("epoch id should be lower or equal than the last id")
		}
		epochIdMap[elem.Id] = true
	}
	// Check for duplicated index in gauge
	gaugeIndexMap := make(map[string]struct{})

	for _, elem := range gs.GaugeList {
		index := string(GaugeKey(elem.PreviousEpochId, elem.PoolId))
		if _, ok := gaugeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for gauge")
		}
		gaugeIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
