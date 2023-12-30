package types

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
