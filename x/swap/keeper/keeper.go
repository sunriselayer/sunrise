package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"

	ibckeeper "github.com/cosmos/ibc-go/v9/modules/core/keeper"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

type Keeper struct {
	appmodule.Environment

	cdc          codec.BinaryCodec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema                  collections.Schema
	Params                  collections.Item[types.Params]
	IncomingInFlightPackets collections.Map[collections.Triple[string, string, uint64], types.IncomingInFlightPacket]
	OutgoingInFlightPackets collections.Map[collections.Triple[string, string, uint64], types.OutgoingInFlightPacket]

	AccountKeeper       types.AccountKeeper
	BankKeeper          types.BankKeeper
	TransferKeeper      types.TransferKeeper
	liquidityPoolKeeper types.LiquidityPoolKeeper

	IbcKeeperFn func() *ibckeeper.Keeper
}

func NewKeeper(
	env appmodule.Environment,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	transferKeeper types.TransferKeeper,
	liquidityPoolKeeper types.LiquidityPoolKeeper,
	ibcKeeperFn func() *ibckeeper.Keeper,
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
		IncomingInFlightPackets: collections.NewMap(
			sb,
			types.IncomingInFlightPacketsKey,
			"incoming_in_flight_packets",
			collections.TripleKeyCodec(collections.StringKey, collections.StringKey, collections.Uint64Key),
			codec.CollValue[types.IncomingInFlightPacket](cdc),
		),
		OutgoingInFlightPackets: collections.NewMap(
			sb,
			types.OutgoingInFlightPacketsKey,
			"outgoing_in_flight_packets",
			collections.TripleKeyCodec(collections.StringKey, collections.StringKey, collections.Uint64Key),
			codec.CollValue[types.OutgoingInFlightPacket](cdc),
		),

		AccountKeeper:       accountKeeper,
		BankKeeper:          bankKeeper,
		TransferKeeper:      transferKeeper,
		liquidityPoolKeeper: liquidityPoolKeeper,
		IbcKeeperFn:         ibcKeeperFn,
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
