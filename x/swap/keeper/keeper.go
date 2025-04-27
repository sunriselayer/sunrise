package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"

	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
	logger       log.Logger
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority string

	addressCodec address.Codec

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
	cdc codec.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	addressCodec address.Codec,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	transferKeeper types.TransferKeeper,
	liquidityPoolKeeper types.LiquidityPoolKeeper,
	ibcKeeperFn func() *ibckeeper.Keeper,
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

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		IncomingInFlightPackets: collections.NewMap(
			sb,
			types.IncomingInFlightPacketsKeyPrefix,
			"incoming_in_flight_packets",
			types.IncomingInFlightPacketsKeyCodec,
			codec.CollValue[types.IncomingInFlightPacket](cdc),
		),
		OutgoingInFlightPackets: collections.NewMap(
			sb,
			types.OutgoingInFlightPacketsKeyPrefix,
			"outgoing_in_flight_packets",
			types.OutgoingInFlightPacketsKeyCodec,
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
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
