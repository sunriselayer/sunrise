package swap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/swap/keeper"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the inFlightPacket
	for _, elem := range genState.OutgoingInFlightPacketList {
		k.SetOutgoingInFlightPacket(ctx, elem)
	}
	// Set all the ackWaitingPacket
	for _, elem := range genState.IncomingInFlightPacketList {
		k.SetIncomingInFlightPacket(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.OutgoingInFlightPacketList = k.GetAllOutgoingInFlightPacket(ctx)
	genesis.IncomingInFlightPacketList = k.GetAllIncomingInFlightPacket(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
