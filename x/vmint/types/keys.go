package types

const (
	// ModuleName defines the module name
	ModuleName = "vmint"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_vmint"

    
)

var (
	ParamsKey = []byte("p_vmint")
)



func KeyPrefix(p string) []byte {
    return []byte(p)
}
