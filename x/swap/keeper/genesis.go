package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/swap/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Set all the incomingPacket
	for _, elem := range genState.IncomingInFlightPackets {
		err := k.SetIncomingInFlightPacket(ctx, elem)
		if err != nil {
			return err
		}
	}
	// Set all the outgoingInFlightPacket
	for _, elem := range genState.OutgoingInFlightPackets {
		err := k.SetOutgoingInFlightPacket(ctx, elem)
		if err != nil {
			return err
		}
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

	genesis.IncomingInFlightPackets, err = k.GetIncomingInFlightPackets(ctx)
	if err != nil {
		return nil, err
	}
	genesis.OutgoingInFlightPackets, err = k.GetOutgoingInFlightPackets(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
