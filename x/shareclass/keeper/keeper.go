package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "cosmossdk.io/x/staking/types"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
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
	Unbondings                *collections.IndexedMap[uint64, types.Unbonding, types.UnbondingsIndexes]
	UnbondingId               collections.Sequence
	RewardMultiplier          collections.Map[collections.Pair[[]byte, string], string]                   // math.Dec
	UsersLastRewardMultiplier collections.Map[collections.Triple[sdk.AccAddress, []byte, string], string] // math.Dec
	LastRewardHandlingTime    collections.Map[[]byte, int64]

	accountKeeper        types.AccountKeeper
	bankKeeper           types.BankKeeper
	stakingKeeper        types.StakingKeeper
	feeKeeper            types.FeeKeeper
	tokenConverterKeeper types.TokenConverterKeeper

	stakingMsgServer stakingtypes.MsgServer
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	feeKeeper types.FeeKeeper,
	tokenConverterKeeper types.TokenConverterKeeper,
	stakingMsgServer stakingtypes.MsgServer,
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
		Unbondings:                collections.NewIndexedMap(sb, types.UnbondingsKeyPrefix, "unbondings", types.UnbondingsKeyCodec, codec.CollValue[types.Unbonding](cdc), types.NewUnbondingsIndexes(sb, addressCodec)),
		UnbondingId:               collections.NewSequence(sb, types.UnbondingIdKey, "unbonding_id"),
		RewardMultiplier:          collections.NewMap(sb, types.RewardMultiplierKeyPrefix, "reward_multiplier", types.RewardMultiplierKeyCodec, collections.StringValue),
		UsersLastRewardMultiplier: collections.NewMap(sb, types.UsersLastRewardMultiplierKeyPrefix, "users_last_reward_multiplier", types.UsersLastRewardMultiplierKeyCodec, collections.StringValue),
		LastRewardHandlingTime:    collections.NewMap(sb, types.LastRewardHandlingTimeKeyPrefix, "last_reward_handling_time", types.LastRewardHandlingTimeKeyCodec, collections.Int64Value),

		accountKeeper:        accountKeeper,
		bankKeeper:           bankKeeper,
		stakingKeeper:        stakingKeeper,
		feeKeeper:            feeKeeper,
		tokenConverterKeeper: tokenConverterKeeper,

		stakingMsgServer: stakingMsgServer,
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
