package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "lending"

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
	// MarketsKey is the prefix for markets
	MarketsKey = collections.NewPrefix("markets")
	// UserPositionsKey is the prefix for user positions
	UserPositionsKey = collections.NewPrefix("user_positions")
	// BorrowsKey is the prefix for borrows
	BorrowsKey = collections.NewPrefix("borrows")
	// BorrowIdKey is the prefix for borrow id sequence
	BorrowIdKey = collections.NewPrefix("borrow_id")
)

// Key codecs
var (
	// UserPositionKeyCodec is the key codec for user positions
	UserPositionKeyCodec = collections.PairKeyCodec(collections.StringKey, collections.StringKey)
)
