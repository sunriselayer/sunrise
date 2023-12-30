package types

const (
	// ModuleName defines the module name
	ModuleName = "grant"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_grant"
)

var (
	ParamsKey = []byte("p_grant")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
