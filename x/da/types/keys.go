package types

import "cosmossdk.io/collections"

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "da"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_da")

// should be changed to use collections
var (
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
