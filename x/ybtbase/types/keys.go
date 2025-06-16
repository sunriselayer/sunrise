package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "ybtbase"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// Collection keys
var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("params")

	// TokensKey is the prefix for tokens collection
	TokensKey = collections.NewPrefix("tokens")

	// GlobalRewardIndexKey is the prefix for global reward index collection
	GlobalRewardIndexKey = collections.NewPrefix("global_reward_index")

	// UserLastRewardIndexKey is the prefix for user last reward index collection
	UserLastRewardIndexKey = collections.NewPrefix("user_last_reward_index")

	// PermissionsKey is the prefix for permissions collection
	PermissionsKey = collections.NewPrefix("permissions")
)
