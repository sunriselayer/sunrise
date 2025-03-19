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
	UnstakingIdsKeyPrefix              = collections.NewPrefix("unstaking_ids/")
	RewardMultiplierKeyPrefix          = collections.NewPrefix("reward_multiplier/")
	UsersLastRewardMultiplierKeyPrefix = collections.NewPrefix("users_last_reward_multiplier/")
)

var (
	UnstakingsKeyCodec                = collections.PairKeyCodec(sdk.AccAddressKey, collections.Uint64Key)
	UnstakingIdsKeyCodec              = sdk.AccAddressKey
	RewardMultiplierKeyCodec          = collections.StringKey
	UsersLastRewardMultiplierKeyCodec = collections.PairKeyCodec(collections.StringKey, sdk.AccAddressKey)
)
