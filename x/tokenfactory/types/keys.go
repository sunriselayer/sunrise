package types

import (
	"strings"

	"cosmossdk.io/collections"
)

const (
	// ModuleName defines the module name
	ModuleName = "tokenfactory"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tokenfactory"

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_tokenfactory")

// KeySeparator is used to combine parts of the keys in the store
const KeySeparator = "|"

var (
	DenomAuthorityMetadataKey      = "authoritymetadata"
	DenomsPrefixKey                = "denoms"
	CreatorPrefixKey               = "creator"
	AdminPrefixKey                 = "admin"
	BeforeSendHookAddressPrefixKey = "beforesendhook"
)

// GetDenomPrefixStore returns the store prefix where all the data associated with a specific denom
// is stored
func GetDenomPrefixStore(denom string) []byte {
	return []byte(strings.Join([]string{DenomsPrefixKey, denom, ""}, KeySeparator))
}

// GetCreatorsPrefix returns the store prefix where the list of the denoms created by a specific
// creator are stored
func GetCreatorPrefix(creator string) []byte {
	return []byte(strings.Join([]string{CreatorPrefixKey, creator, ""}, KeySeparator))
}

// GetCreatorsPrefix returns the store prefix where a list of all creator addresses are stored
func GetCreatorsPrefix() []byte {
	return []byte(strings.Join([]string{CreatorPrefixKey, ""}, KeySeparator))
}
