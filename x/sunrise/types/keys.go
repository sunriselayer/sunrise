package types

const (
	// ModuleName defines the module name
	ModuleName = "sunrise"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_sunrise"
)

var (
	ParamsKey = []byte("p_sunrise")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
