package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
	logger       log.Logger

	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority string

	addressCodec address.Codec

	Schema     collections.Schema
	Params     collections.Item[types.Params]
	Pools      collections.Map[uint64, types.Pool]
	PoolId     collections.Sequence
	Positions  *collections.IndexedMap[uint64, types.Position, types.PositionsIndexes]
	PositionId collections.Sequence

	bankKeeper types.BankKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	addressCodec address.Codec,
	bankKeeper types.BankKeeper,
) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		logger:       logger,
		authority:    authority,
		addressCodec: addressCodec,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Pools:  collections.NewMap(sb, types.PoolsKeyPrefix, "pools", types.PoolsKeyCodec, codec.CollValue[types.Pool](cdc)),
		PoolId: collections.NewSequence(sb, types.PoolIdKey, "pool_id"),
		Positions: collections.NewIndexedMap(
			sb,
			types.PositionsKeyPrefix,
			"positions",
			types.PositionsKeyCodec,
			codec.CollValue[types.Position](cdc),
			types.NewPositionsIndexes(sb, addressCodec),
		),
		PositionId: collections.NewSequence(sb, types.PositionIdKey, "position_id"),

		bankKeeper: bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
