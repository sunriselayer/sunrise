package types

const (
	// ModuleName defines the module name
	ModuleName = "feeconverter"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_feeconverter"
)

var (
	ParamsKey = []byte("p_feeconverter")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
