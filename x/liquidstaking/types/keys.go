package types

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "liquidstaking"

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

	UnstakingsKeyPrefix                = collections.NewPrefix("unstakings/")
	UnstakingIdKey                     = collections.NewPrefix("unstaking_id/")
	RewardMultiplierKeyPrefix          = collections.NewPrefix("reward_multiplier/")
	UsersLastRewardMultiplierKeyPrefix = collections.NewPrefix("users_last_reward_multiplier/")
	LastRewardHandlingTimeKeyPrefix    = collections.NewPrefix("last_reward_handling_time/")
)

var (
	UnstakingsKeyCodec                = collections.PairKeyCodec(collections.Int64Key, collections.Uint64Key)
	RewardMultiplierKeyCodec          = collections.PairKeyCodec(collections.BytesKey, collections.StringKey)
	UsersLastRewardMultiplierKeyCodec = collections.TripleKeyCodec(sdk.AccAddressKey, collections.BytesKey, collections.StringKey)
	LastRewardHandlingTimeKeyCodec    = collections.BytesKey
)
