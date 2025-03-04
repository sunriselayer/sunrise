package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Set all the epoch
	for _, elem := range genState.Epochs {
		k.SetEpoch(ctx, elem)
	}

	// Set epoch count
	k.SetEpochCount(ctx, genState.EpochCount)
	// Set all the gauges
	for _, elem := range genState.Gauges {
		k.SetGauge(ctx, elem)
	}
	// Set all the votes
	for _, vote := range genState.Votes {
		k.SetVote(ctx, vote)
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

	genesis.Epochs = k.GetAllEpoch(ctx)
	genesis.EpochCount = k.GetEpochCount(ctx)
	genesis.Gauges = k.GetAllGauges(ctx)
	genesis.Votes = k.GetAllVotes(ctx)

	return genesis, nil
}
