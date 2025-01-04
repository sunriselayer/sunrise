package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/swap/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Set all the incomingPacket
	for _, elem := range genState.IncomingInFlightPackets {
		k.SetIncomingInFlightPacket(ctx, elem)
	}
	// Set all the outgoingInFlightPacket
	for _, elem := range genState.OutgoingInFlightPackets {
		k.SetOutgoingInFlightPacket(ctx, elem)
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	genesis.IncomingInFlightPackets = k.GetIncomingInFlightPackets(ctx)
	genesis.OutgoingInFlightPackets = k.GetOutgoingInFlightPackets(ctx)

	return genesis, nil
}
