package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PoolList: []Pool{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in pool
	poolIdMap := make(map[uint64]bool)
	poolCount := gs.GetPoolCount()
	for _, elem := range gs.PoolList {
		if _, ok := poolIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for pool")
		}
		if elem.Id >= poolCount {
			return fmt.Errorf("pool id should be lower or equal than the last id")
		}
		poolIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
