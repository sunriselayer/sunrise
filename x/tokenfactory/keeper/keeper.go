package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/tokenfactory/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema            collections.Schema
	Params            collections.Item[types.Params]
	AuthorityMetadata collections.Map[string, types.DenomAuthorityMetadata]
	CreatorAddresses  collections.Map[string, []byte]
	// BeforeSendHook    collections.Map[string, []byte]
	DenomFromCreator collections.Map[collections.Pair[sdk.AccAddress, string], []byte]

	accountKeeper      types.AccountKeeper
	bankKeeper         types.BankKeeper
	distributionKeeper types.DistributionKeeper
	// contractKeeper     types.ContractKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distributionKeeper types.DistributionKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params:            collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		AuthorityMetadata: collections.NewMap(sb, types.DenomAuthorityMetadataKey, "authority_metadata", collections.StringKey, codec.CollValue[types.DenomAuthorityMetadata](cdc)),
		CreatorAddresses:  collections.NewMap(sb, types.CreatorsKeyPrefix, "creator_addresses", collections.StringKey, collections.BytesValue),
		// BeforeSendHook:    collections.NewMap(sb, types.BeforeSendHookAddressPrefixKey, "before_send_hook", collections.StringKey, collections.BytesValue),
		DenomFromCreator: collections.NewMap(sb, types.DenomFromCreatorKey, "denom_from_creator", collections.PairKeyCodec(sdk.AccAddressKey, collections.StringKey), collections.BytesValue),

		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		distributionKeeper: distributionKeeper,
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

// // Set the wasm keeper.
// func (k *Keeper) SetContractKeeper(contractKeeper types.ContractKeeper) {
// 	k.contractKeeper = contractKeeper
// }
