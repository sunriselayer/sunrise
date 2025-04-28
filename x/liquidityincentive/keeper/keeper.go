package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

type Keeper struct {
	appmodule.Environment

	cdc          codec.BinaryCodec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema              collections.Schema
	Params              collections.Item[types.Params]
	Epochs              collections.Map[uint64, types.Epoch]
	EpochId             collections.Sequence
	Gauges              collections.Map[collections.Pair[uint64, uint64], types.Gauge]
	Votes               collections.Map[sdk.AccAddress, types.Vote]
	Bribes              *collections.IndexedMap[uint64, types.Bribe, types.BribesIndexes]
	BribeId             collections.Sequence
	BribeAllocations    collections.Map[collections.Triple[sdk.AccAddress, uint64, uint64], types.BribeAllocation]
	BribeExpiredEpochId collections.Item[uint64]

	accountKeeper        types.AccountKeeper
	bankKeeper           types.BankKeeper
	stakingKeeper        types.StakingKeeper
	feeKeeper            types.FeeKeeper
	tokenConverterKeeper types.TokenConverterKeeper
	liquidityPoolKeeper  types.LiquidityPoolKeeper
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	feeKeeper types.FeeKeeper,
	tokenConverterKeeper types.TokenConverterKeeper,
	liquidityPoolKeeper types.LiquidityPoolKeeper,
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

		Params:  collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Epochs:  collections.NewMap(sb, types.EpochsKeyPrefix, "epochs", types.EpochsKeyCodec, codec.CollValue[types.Epoch](cdc)),
		EpochId: collections.NewSequence(sb, types.EpochIdKey, "epoch_id"),
		Gauges:  collections.NewMap(sb, types.GaugesKeyPrefix, "gauges", types.GaugesKeyCodec, codec.CollValue[types.Gauge](cdc)),
		Votes:   collections.NewMap(sb, types.VotesKeyPrefix, "votes", types.VotesKeyCodec, codec.CollValue[types.Vote](cdc)),
		Bribes: collections.NewIndexedMap(
			sb,
			types.BribesKeyPrefix,
			"bribes",
			types.BribesKeyCodec,
			codec.CollValue[types.Bribe](cdc),
			types.NewBribesIndexes(sb),
		),
		BribeId: collections.NewSequence(sb, types.BribeIdKey, "bribe_id"),
		BribeAllocations: collections.NewMap(
			sb,
			types.BribeAllocationsKeyPrefix,
			"bribe_allocations",
			types.BribeAllocationsKeyCodec,
			codec.CollValue[types.BribeAllocation](cdc),
		),
		BribeExpiredEpochId: collections.NewItem(sb, types.BribeExpiredEpochIdKey, "bribe_expired_epoch_id", collections.Uint64Value),

		accountKeeper:        authKeeper,
		bankKeeper:           bankKeeper,
		stakingKeeper:        stakingKeeper,
		feeKeeper:            feeKeeper,
		tokenConverterKeeper: tokenConverterKeeper,
		liquidityPoolKeeper:  liquidityPoolKeeper,
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
