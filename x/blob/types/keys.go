package types

const (
	// ModuleName defines the module name
	ModuleName = "blob"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_blob"
)

var (
	ParamsKey = []byte("p_blob")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
