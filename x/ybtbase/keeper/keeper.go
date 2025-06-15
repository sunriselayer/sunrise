package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	authKeeper types.AuthKeeper
	bankKeeper types.BankKeeper

	Schema collections.Schema
	Params collections.Item[types.Params]

	// Collections for YBT base tokens
	Tokens              collections.Map[string, types.Token]
	GlobalRewardIndex   collections.Map[string, math.LegacyDec]
	UserLastRewardIndex collections.Map[collections.Pair[string, string], math.LegacyDec]
	YieldPermissions    collections.Map[collections.Pair[string, string], bool]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		authKeeper:   authKeeper,
		bankKeeper:   bankKeeper,

		Params:              collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Tokens:              collections.NewMap(sb, collections.NewPrefix(0), "tokens", collections.StringKey, codec.CollValue[types.Token](cdc)),
		GlobalRewardIndex:   collections.NewMap(sb, collections.NewPrefix(1), "global_reward_index", collections.StringKey, sdk.LegacyDecValue),
		UserLastRewardIndex: collections.NewMap(sb, collections.NewPrefix(2), "user_last_reward_index", collections.PairKeyCodec(collections.StringKey, collections.StringKey), sdk.LegacyDecValue),
		YieldPermissions:    collections.NewMap(sb, collections.NewPrefix(3), "yield_permissions", collections.PairKeyCodec(collections.StringKey, collections.StringKey), collections.BoolValue),
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
