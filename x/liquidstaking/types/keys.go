package types

const (
	// ModuleName defines the module name
	ModuleName = "liquidstaking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_liquidstaking"
)

var (
	ParamsKey = []byte("p_liquidstaking")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
