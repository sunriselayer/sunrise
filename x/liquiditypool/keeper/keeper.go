package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

type Keeper struct {
	appmodule.Environment

	cdc          codec.BinaryCodec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema     collections.Schema
	Params     collections.Item[types.Params]
	Pools      collections.Map[uint64, types.Pool]
	PoolId     collections.Sequence
	Positions  *collections.IndexedMap[uint64, types.Position, types.PositionsIndexes]
	PositionId collections.Sequence

	bankKeeper types.BankKeeper
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	bankKeeper types.BankKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(env.KVStoreService)

	k := Keeper{
		Environment:  env,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

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
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

func (k Keeper) GetNextPoolID(ctx context.Context) (uint64, error) {
	count, err := k.GetPoolCount(ctx)
	if err != nil {
		return 0, err
	}
	return count + 1, nil
}
