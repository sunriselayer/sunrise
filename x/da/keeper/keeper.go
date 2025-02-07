package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/sunriselayer/sunrise/x/da/types"
)

type Keeper struct {
	appmodule.Environment

	cdc          codec.BinaryCodec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema          collections.Schema
	Params          collections.Item[types.Params]
	PublishedData   *collections.IndexedMap[string, types.PublishedData, types.PublishedDataIndexes]
	ChallengeCounts collections.Item[uint64]
	FaultCounts     collections.Map[[]byte, uint64]
	Proofs          collections.Map[collections.Pair[string, []byte], types.Proof]
	Invalidities    collections.Map[collections.Pair[string, []byte], types.Invalidity]

	BankKeeper     types.BankKeeper
	StakingKeeper  types.StakingKeeper
	SlashingKeeper types.SlashingKeeper
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	slashingKeeper types.SlashingKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(env.KVStoreService)

	k := Keeper{
		Environment:  env,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		PublishedData: collections.NewIndexedMap(
			sb,
			types.PublishedDataKeyPrefix,
			"published_data",
			types.PublishedDataKeyCodec,
			codec.CollValue[types.PublishedData](cdc),
			types.NewPublishedDataIndexes(sb),
		),
		ChallengeCounts: collections.NewItem(sb, types.ChallengeCountsKeyPrefix, "challenge_counts", collections.Uint64Value),
		FaultCounts:     collections.NewMap(sb, types.FaultCountsKeyPrefix, "fault_counts", types.FaultCounterKeyCodec, collections.Uint64Value),
		Proofs:          collections.NewMap(sb, types.ProofKeyPrefix, "proofs", types.ProofKeyCodec, codec.CollValue[types.Proof](cdc)),
		Invalidities:    collections.NewMap(sb, types.InvalidityKeyPrefix, "invalidities", types.ProofKeyCodec, codec.CollValue[types.Invalidity](cdc)),

		BankKeeper:     bankKeeper,
		StakingKeeper:  stakingKeeper,
		SlashingKeeper: slashingKeeper,
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
