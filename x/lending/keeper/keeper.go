package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/sunriselayer/sunrise/x/lending/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService corestore.KVStoreService
	logger       log.Logger

	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority string

	addressCodec address.Codec

	Schema collections.Schema
	Params collections.Item[types.Params]
	
	// Markets stores all lending markets by denom
	Markets collections.Map[string, types.Market]
	// UserPositions stores user positions by (user_address, denom)
	UserPositions collections.Map[collections.Pair[string, string], types.UserPosition]
	// Borrows stores all borrow positions by id
	Borrows collections.Map[uint64, types.Borrow]
	// BorrowId tracks the next borrow id
	BorrowId collections.Sequence

	bankKeeper types.BankKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService corestore.KVStoreService,
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

		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Markets:       collections.NewMap(sb, types.MarketsKey, "markets", collections.StringKey, codec.CollValue[types.Market](cdc)),
		UserPositions: collections.NewMap(sb, types.UserPositionsKey, "user_positions", collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[types.UserPosition](cdc)),
		Borrows:       collections.NewMap(sb, types.BorrowsKey, "borrows", collections.Uint64Key, codec.CollValue[types.Borrow](cdc)),
		BorrowId:      collections.NewSequence(sb, types.BorrowIdKey, "borrow_id"),

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
