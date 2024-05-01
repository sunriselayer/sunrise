package types

import (
"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PairList: []Pair{},
PoolList: []Pool{},
TwapList: []Twap{},
// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in pair
pairIndexMap := make(map[string]struct{})

for _, elem := range gs.PairList {
	index := string(PairKey(elem.Index))
	if _, ok := pairIndexMap[index]; ok {
		return fmt.Errorf("duplicated index for pair")
	}
	pairIndexMap[index] = struct{}{}
}
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
// Check for duplicated index in twap
twapIndexMap := make(map[string]struct{})

for _, elem := range gs.TwapList {
	index := string(TwapKey(elem.Index))
	if _, ok := twapIndexMap[index]; ok {
		return fmt.Errorf("duplicated index for twap")
	}
	twapIndexMap[index] = struct{}{}
}
// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
