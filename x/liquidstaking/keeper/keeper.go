package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

type Keeper struct {
	appmodule.Environment

	cdc          codec.BinaryCodec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema                    collections.Schema
	Params                    collections.Item[types.Params]
	Unstakings                collections.Map[collections.Pair[sdk.AccAddress, uint64], types.Unstaking]
	UnstakingIds              collections.Map[sdk.AccAddress, uint64]
	RewardMultiplier          collections.Map[string, string]
	UsersLastRewardMultiplier collections.Map[collections.Pair[string, sdk.AccAddress], string]

	accountKeeper        types.AccountKeeper
	bankKeeper           types.BankKeeper
	stakingKeeper        types.StakingKeeper
	tokenConverterKeeper types.TokenConverterKeeper
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	tokenConverterKeeper types.TokenConverterKeeper,
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

		Params:                    collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Unstakings:                collections.NewMap(sb, types.UnstakingsKeyPrefix, "unstakings", types.UnstakingsKeyCodec, codec.CollValue[types.Unstaking](cdc)),
		UnstakingIds:              collections.NewMap(sb, types.UnstakingIdsKeyPrefix, "unstaking_ids", types.UnstakingIdsKeyCodec, collections.Uint64Value),
		RewardMultiplier:          collections.NewMap(sb, types.RewardMultiplierKeyPrefix, "reward_multiplier", types.RewardMultiplierKeyCodec, collections.StringValue),
		UsersLastRewardMultiplier: collections.NewMap(sb, types.UsersLastRewardMultiplierKeyPrefix, "users_last_reward_multiplier", types.UsersLastRewardMultiplierKeyCodec, collections.StringValue),

		accountKeeper:        accountKeeper,
		bankKeeper:           bankKeeper,
		stakingKeeper:        stakingKeeper,
		tokenConverterKeeper: tokenConverterKeeper,
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
