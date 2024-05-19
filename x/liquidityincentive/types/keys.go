package types

const (
	// ModuleName defines the module name
	ModuleName = "liquidityincentive"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_liquidityincentive"
)

var (
	ParamsKey = []byte("p_liquidityincentive")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
