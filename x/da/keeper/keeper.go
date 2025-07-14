package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/sunriselayer/sunrise/x/da/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
	logger       log.Logger

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	addressCodec address.Codec

	Schema               collections.Schema
	Params               collections.Item[types.Params]
	PublishedData        *collections.IndexedMap[string, types.PublishedData, types.PublishedDataIndexes]
	ChallengeCounts      collections.Item[uint64]
	FaultCounts          collections.Map[[]byte, uint64]
	Proofs               collections.Map[collections.Pair[string, []byte], types.Proof]
	Invalidities         collections.Map[collections.Pair[string, []byte], types.Invalidity]
	ProofDeputies        collections.Map[[]byte, []byte]
	LastSlashBlockHeight collections.Item[int64]

	BankKeeper     types.BankKeeper
	StakingKeeper  types.StakingKeeper
	SlashingKeeper types.SlashingKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	addressCodec address.Codec,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	slashingKeeper types.SlashingKeeper,
) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		logger:       logger,
		authority:    authority,
		addressCodec: addressCodec,

		Params:               collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		PublishedData:        collections.NewIndexedMap(sb, types.PublishedDataKeyPrefix, "published_data", types.PublishedDataKeyCodec, codec.CollValue[types.PublishedData](cdc), types.NewPublishedDataIndexes(sb)),
		ChallengeCounts:      collections.NewItem(sb, types.ChallengeCountsKeyPrefix, "challenge_counts", collections.Uint64Value),
		FaultCounts:          collections.NewMap(sb, types.FaultCountsKeyPrefix, "fault_counts", types.FaultCounterKeyCodec, collections.Uint64Value),
		Proofs:               collections.NewMap(sb, types.ProofKeyPrefix, "proofs", types.ProofKeyCodec, codec.CollValue[types.Proof](cdc)),
		Invalidities:         collections.NewMap(sb, types.InvalidityKeyPrefix, "invalidities", types.InvalidityKeyCodec, codec.CollValue[types.Invalidity](cdc)),
		ProofDeputies:        collections.NewMap(sb, types.ProofDeputiesKeyPrefix, "proof_deputy", types.ProofDeputyKeyCodec, collections.BytesValue),
		LastSlashBlockHeight: collections.NewItem(sb, types.LastSlashBlockHeightKey, "last_slash_block_height", collections.Int64Value),

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
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
