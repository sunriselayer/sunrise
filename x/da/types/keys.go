package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "da"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_da"
)

var (
	ParamsKey = []byte("p_da")

	PublishedDataKeyPrefix = []byte("published_data/")
	FaultCounterKeyPrefix  = []byte("fault_counter/")
	ProofKeyPrefix         = []byte("proof/")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func PublishedDataKey(metadataUri string) []byte {
	return append(PublishedDataKeyPrefix, metadataUri...)
}

func GetFaultCounterKey(val sdk.ValAddress) []byte {
	return append(FaultCounterKeyPrefix, address.MustLengthPrefix(val)...)
}

func ProofKey(metadataUri string, sender string) []byte {
	return append(append(ProofKeyPrefix, metadataUri...), sender...)
}
