package types

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "shareclass"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("params/")

	UnbondingsKeyPrefix                 = collections.NewPrefix("unbondings/")
	UnbondingsAddressIndexPrefix        = collections.NewPrefix("unbondings_by_address/")
	UnbondingsCompletionTimeIndexPrefix = collections.NewPrefix("unbondings_by_completion_time/")
	UnbondingIdKey                      = collections.NewPrefix("unbonding_id/")
	RewardMultiplierKeyPrefix           = collections.NewPrefix("reward_multiplier/")
	UsersLastRewardMultiplierKeyPrefix  = collections.NewPrefix("users_last_reward_multiplier/")
	LastRewardHandlingTimeKeyPrefix     = collections.NewPrefix("last_reward_handling_time/")
)

type UnbondingsIndexes struct {
	Address        *indexes.Multi[sdk.AccAddress, uint64, Unbonding]
	CompletionTime *indexes.Multi[int64, uint64, Unbonding]
}

func (i UnbondingsIndexes) IndexesList() []collections.Index[uint64, Unbonding] {
	return []collections.Index[uint64, Unbonding]{
		i.Address,
		i.CompletionTime,
	}
}

func NewUnbondingsIndexes(sb *collections.SchemaBuilder, addressCodec address.Codec) UnbondingsIndexes {
	return UnbondingsIndexes{
		Address: indexes.NewMulti(
			sb,
			UnbondingsAddressIndexPrefix,
			"unbondings_by_address",
			sdk.AccAddressKey,
			collections.Uint64Key,
			func(_ uint64, v Unbonding) (sdk.AccAddress, error) {
				return addressCodec.StringToBytes(v.Address)
			},
		),
		CompletionTime: indexes.NewMulti(
			sb,
			UnbondingsCompletionTimeIndexPrefix,
			"unbondings_by_completion_time",
			collections.Int64Key,
			collections.Uint64Key,
			func(_ uint64, v Unbonding) (int64, error) {
				return v.CompletionTime.Unix(), nil
			},
		),
	}
}

var (
	UnbondingsKeyCodec                = collections.Uint64Key
	RewardMultiplierKeyCodec          = collections.PairKeyCodec(collections.BytesKey, collections.StringKey)
	UsersLastRewardMultiplierKeyCodec = collections.TripleKeyCodec(sdk.AccAddressKey, collections.BytesKey, collections.StringKey)
	LastRewardHandlingTimeKeyCodec    = collections.BytesKey
)
