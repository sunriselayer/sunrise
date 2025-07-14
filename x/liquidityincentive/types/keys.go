package types

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "liquidityincentive"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_liquidityincentive"

	// BribeAccount is the account that holds the bribe funds
	BribeAccount = ModuleName + "/bribe"

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("params/")

	EpochsKeyPrefix                 = collections.NewPrefix("epochs/")
	EpochIdKey                      = collections.NewPrefix("epoch_id/")
	VotesKeyPrefix                  = collections.NewPrefix("votes/")
	BribesKeyPrefix                 = collections.NewPrefix("bribes/")
	BribeIdKey                      = collections.NewPrefix("bribe_id/")
	BribesEpochIdIndexPrefix        = collections.NewPrefix("bribes_by_epoch_id/")
	BribesPoolIdIndexPrefix         = collections.NewPrefix("bribes_by_pool_id/")
	BribeAllocationsKeyPrefix       = collections.NewPrefix("bribe_allocations/")
	BribeAllocationsByEpochIdPrefix = collections.NewPrefix("bribe_allocations_by_epoch_id/")
	BribeExpiredEpochIdKey          = collections.NewPrefix("bribe_expired_epoch_id/")
)

type BribesIndexes struct {
	EpochId *indexes.Multi[uint64, uint64, Bribe]
	PoolId  *indexes.Multi[uint64, uint64, Bribe]
}

func (i BribesIndexes) IndexesList() []collections.Index[uint64, Bribe] {
	return []collections.Index[uint64, Bribe]{
		i.EpochId,
		i.PoolId,
	}
}

func NewBribesIndexes(sb *collections.SchemaBuilder) BribesIndexes {
	return BribesIndexes{
		EpochId: indexes.NewMulti(sb,
			BribesEpochIdIndexPrefix,
			"bribes_by_epoch_id",
			collections.Uint64Key,
			collections.Uint64Key,
			func(_ uint64, v Bribe) (uint64, error) {
				return v.EpochId, nil
			},
		),
		PoolId: indexes.NewMulti(sb,
			BribesPoolIdIndexPrefix,
			"bribes_by_pool_id",
			collections.Uint64Key,
			collections.Uint64Key,
			func(_ uint64, v Bribe) (uint64, error) {
				return v.PoolId, nil
			},
		),
	}
}

type BribeAllocationsIndexes struct {
	EpochId *indexes.Multi[uint64, collections.Triple[sdk.AccAddress, uint64, uint64], BribeAllocation]
}

func (i BribeAllocationsIndexes) IndexesList() []collections.Index[collections.Triple[sdk.AccAddress, uint64, uint64], BribeAllocation] {
	return []collections.Index[collections.Triple[sdk.AccAddress, uint64, uint64], BribeAllocation]{i.EpochId}
}

func NewBribeAllocationsIndexes(sb *collections.SchemaBuilder) BribeAllocationsIndexes {
	return BribeAllocationsIndexes{
		EpochId: indexes.NewMulti(sb,
			BribeAllocationsByEpochIdPrefix,
			"bribe_allocations_by_epoch_id",
			collections.Uint64Key,
			BribeAllocationsKeyCodec,
			func(pk collections.Triple[sdk.AccAddress, uint64, uint64], v BribeAllocation) (uint64, error) {
				return v.EpochId, nil
			},
		),
	}
}

var (
	EpochsKeyCodec           = collections.Uint64Key
	VotesKeyCodec            = sdk.AccAddressKey
	BribesKeyCodec           = collections.Uint64Key
	BribeAllocationsKeyCodec = collections.TripleKeyCodec(sdk.AccAddressKey, collections.Uint64Key, collections.Uint64Key)
)

func GaugeKey(previousEpochId uint64, poolId uint64) collections.Pair[uint64, uint64] {
	return collections.Join(previousEpochId, poolId)
}

func BribeKey(epochId uint64, poolId uint64) collections.Pair[uint64, uint64] {
	return collections.Join(epochId, poolId)
}

func BribeAllocationKey(address sdk.AccAddress, epochId uint64, poolId uint64) collections.Triple[sdk.AccAddress, uint64, uint64] {
	return collections.Join3(address, epochId, poolId)
}
