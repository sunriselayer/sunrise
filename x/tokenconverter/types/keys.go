package types

const (
	// ModuleName defines the module name
	ModuleName = "tokenconverter"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tokenconverter"
)

var (
	ParamsKey = []byte("p_tokenconverter")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
