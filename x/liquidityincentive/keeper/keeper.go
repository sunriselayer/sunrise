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

	Schema          collections.Schema
	Params          collections.Item[types.Params]
	Epochs          collections.Map[uint64, types.Epoch]
	EpochId         collections.Sequence
	Gauges          collections.Map[collections.Pair[uint64, uint64], types.Gauge]
	Votes           collections.Map[sdk.AccAddress, types.Vote]
	Bribes          collections.Map[collections.Pair[uint64, uint64], types.Bribe]
	UnclaimedBribes collections.Map[collections.Triple[sdk.AccAddress, uint64, uint64], types.UnclaimedBribe]

	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	stakingKeeper       types.StakingKeeper
	liquidityPoolKeeper types.LiquidityPoolKeeper
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
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
		Bribes:  collections.NewMap(sb, types.BribesKeyPrefix, "bribes", types.BribesKeyCodec, codec.CollValue[types.Bribe](cdc)),
		UnclaimedBribes: collections.NewMap(
			sb,
			types.UnclaimedBribesKeyPrefix,
			"unclaimed_bribes",
			types.UnclaimedBribesKeyCodec,
			codec.CollValue[types.UnclaimedBribe](cdc),
		),

		accountKeeper:       authKeeper,
		bankKeeper:          bankKeeper,
		stakingKeeper:       stakingKeeper,
		liquidityPoolKeeper: liquidityPoolKeeper,
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
