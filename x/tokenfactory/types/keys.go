package types

import (
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

var (
	// ParamsKey is the prefix for tokenfactory params
	ParamsKey = collections.NewPrefix("p_tokenfactory")
	// DenomAuthorityMetadataKey is the prefix for the DenomAuthorityMetadata map.
	DenomAuthorityMetadataKey = collections.NewPrefix("authority_metadata/")
	// CreatorsKeyPrefix is the prefix for the CreatorAddresses map.
	CreatorsKeyPrefix = collections.NewPrefix("creator_addresses/")
	// DenomFromCreatorKey is the prefix for the DenomFromCreator map.
	DenomFromCreatorKey = collections.NewPrefix("denom_from_creator/")
)
