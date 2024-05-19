package types

const (
	// ModuleName defines the module name
	ModuleName = "liquiditypool"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_liquiditypool"
)

var (
	ParamsKey = []byte("p_liquiditypool")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	PoolKey      = "Pool/value/"
	PoolCountKey = "Pool/count/"
)

const (
	PositionKey      = "Position/value/"
	PositionCountKey = "Position/count/"
)
