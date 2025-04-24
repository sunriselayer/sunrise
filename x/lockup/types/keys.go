package types

import (
	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "lockup"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

var (
	ParamsKey                    = collections.NewPrefix("params")
	NextLockupAccountIdKeyPrefix = collections.NewPrefix("next_lockup_account_id")
	LockupAccountsKeyPrefix      = collections.NewPrefix("lockup_accounts")
)

var (
	NextLockupAccountIdKeyCodec = sdk.AccAddressKey
	LockupAccountsKeyCodec      = collections.PairKeyCodec(collections.BytesKey, collections.Uint64Key)
)
