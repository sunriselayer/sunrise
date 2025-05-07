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

	Schema                   collections.Schema
	Params                   collections.Item[types.Params]
	CommitmentKeys           collections.Map[[]byte, types.CommitmentKey]
	BlobDeclarations         *collections.IndexedMap[[]byte, types.BlobDeclaration, types.BlobDeclarationIndexes]
	ValidatorsPowerSnapshots collections.Map[int64, types.ValidatorsPowerSnapshot]
	BlobCommitments          *collections.IndexedMap[[]byte, types.BlobCommitment, types.BlobCommitmentIndexes]
	Challenges               *collections.IndexedMap[uint64, types.Challenge, types.ChallengeIndexes]
	ChallengeId              collections.Sequence

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

		Params:                   collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		CommitmentKeys:           collections.NewMap(sb, types.CommitmentKeysKeyPrefix, "commitment_keys", types.CommitmentKeyCodec, codec.CollValue[types.CommitmentKey](cdc)),
		BlobDeclarations:         collections.NewIndexedMap(sb, types.BlobDeclarationsKeyPrefix, "blob_declarations", types.BlobDeclarationKeyCodec, codec.CollValue[types.BlobDeclaration](cdc), types.NewBlobDeclarationIndexes(sb)),
		ValidatorsPowerSnapshots: collections.NewMap(sb, types.ValidatorsPowerSnapshotsKeyPrefix, "validators_power_snapshots", types.ValidatorsPowerSnapshotKeyCodec, codec.CollValue[types.ValidatorsPowerSnapshot](cdc)),
		BlobCommitments:          collections.NewIndexedMap(sb, types.BlobCommitmentsKeyPrefix, "blob_commitments", types.BlobCommitmentKeyCodec, codec.CollValue[types.BlobCommitment](cdc), types.NewBlobCommitmentIndexes(sb)),
		Challenges:               collections.NewIndexedMap(sb, types.ChallengesKeyPrefix, "challenges", types.ChallengeKeyCodec, codec.CollValue[types.Challenge](cdc), types.NewChallengeIndexes(sb)),
		ChallengeId:              collections.NewSequence(sb, types.ChallengeIdKeyPrefix, "challenge_id"),

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
