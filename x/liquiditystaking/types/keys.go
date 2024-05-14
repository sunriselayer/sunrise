package types

const (
	// ModuleName defines the module name
	ModuleName = "liquiditystaking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_liquiditystaking"
)

var (
	ParamsKey = []byte("p_liquiditystaking")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
