package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "stream"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_stream"
)

var (
	ParamsKey = []byte("p_stream")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	// AttestationRequestKey indexes attestation requests by nonce
	AttestationRequestKey = "AttestationRequestKey"

	// LatestUnBondingBlockHeight indexes the latest validator unbonding block
	// height
	LatestUnBondingBlockHeight = "LatestUnBondingBlockHeight"

	// LatestAttestationNonce indexes the latest attestation request nonce
	LatestAttestationNonce = "LatestAttestationNonce"

	// EarliestAvailableAttestationNonce indexes the earliest available
	// attestation nonce
	EarliestAvailableAttestationNonce = "EarliestAvailableAttestationNonce"

	// EVMAddress indexes evm addresses by validator address
	EvmAddress = "EvmAddress"
)

// GetAttestationKey returns the following key format
// prefix    nonce
// [0x0][0 0 0 0 0 0 0 1]
func GetAttestationKey(nonce uint64) string {
	return AttestationRequestKey + string(UInt64Bytes(nonce))
}

func ConvertByteArrToString(value []byte) string {
	var ret strings.Builder
	for i := 0; i < len(value); i++ {
		ret.WriteString(string(value[i]))
	}
	return ret.String()
}

func GetEVMKey(valAddress sdk.ValAddress) []byte {
	return append([]byte(EvmAddress), valAddress...)
}
