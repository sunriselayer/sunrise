package types

import (
	"cosmossdk.io/collections"
)

const (
	// ModuleName defines the module name
	ModuleName = "selfdelegation"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey                      = collections.NewPrefix("params/")
	LockupAccountsKeyPrefix        = collections.NewPrefix("lockup_accounts/")
	SelfDelegationProxiesKeyPrefix = collections.NewPrefix("self_delegation_proxies/")
)

var (
	LockupAccountsKeyCodec        = collections.BytesKey
	SelfDelegationProxiesKeyCodec = collections.BytesKey
)
