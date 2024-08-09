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

	PublishedDataKeyPrefix        = []byte("published_data/")
	UnverifiedDataByTimeKeyPrefix = []byte("unverified_data_by_time/")
	FaultCounterKeyPrefix         = []byte("fault_counter/")
	ProofKeyPrefix                = []byte("proof/")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func PublishedDataKey(metadataUri string) []byte {
	return append(PublishedDataKeyPrefix, metadataUri...)
}

func UnverifiedDataTimeKeyPrefix(timestamp uint64) []byte {
	return append(UnverifiedDataByTimeKeyPrefix, sdk.Uint64ToBigEndian(timestamp)...)
}

func UnverifiedDataByTimeKey(timestamp uint64, metadataUri string) []byte {
	return append(UnverifiedDataTimeKeyPrefix(timestamp), metadataUri...)
}

func GetFaultCounterKey(val sdk.ValAddress) []byte {
	return append(FaultCounterKeyPrefix, address.MustLengthPrefix(val)...)
}

func ProofKey(metadataUri string, sender string) []byte {
	return append(append(ProofKeyPrefix, metadataUri...), sender...)
}
