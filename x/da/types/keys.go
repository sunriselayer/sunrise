package types

const (
	// ModuleName defines the module name
	ModuleName = "da"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_da"
)

var (
	ParamsKey = []byte("p_da")

	PublishedDataKeyPrefix = []byte("published_data/")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func PublishedDataKey(dataHash []byte) []byte {
	return append(PublishedDataKeyPrefix, dataHash...)
}
