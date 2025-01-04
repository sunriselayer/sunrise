package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/sunriselayer/sunrise/x/da/types"

	// should be removed
	"cosmossdk.io/core/store"
)

type Keeper struct {
	appmodule.Environment

	cdc          codec.BinaryCodec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	BankKeeper     types.BankKeeper
	StakingKeeper  types.StakingKeeper
	SlashingKeeper types.SlashingKeeper

	// should be removed
	storeService store.KVStoreService
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	slashingKeeper types.SlashingKeeper,
	// should be removed
	storeService store.KVStoreService,
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

		BankKeeper:     bankKeeper,
		StakingKeeper:  stakingKeeper,
		SlashingKeeper: slashingKeeper,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),

		// should be removed
		storeService: storeService,
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
