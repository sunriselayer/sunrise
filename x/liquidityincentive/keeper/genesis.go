package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Set all the epoch
	for _, elem := range genState.Epochs {
		err := k.SetEpoch(ctx, elem)
		if err != nil {
			return err
		}
	}

	// Set epoch count
	err := k.SetEpochCount(ctx, genState.EpochCount)
	if err != nil {
		return err
	}
	// Set all the gauges
	for _, elem := range genState.Gauges {
		err := k.SetGauge(ctx, elem)
		if err != nil {
			return err
		}
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

	genesis.Epochs, err = k.GetAllEpoch(ctx)
	if err != nil {
		return nil, err
	}
	genesis.EpochCount, err = k.GetEpochCount(ctx)
	if err != nil {
		return nil, err
	}
	genesis.Gauges, err = k.GetAllGauges(ctx)
	if err != nil {
		return nil, err
	}
	genesis.Votes = k.GetAllVotes(ctx)

	return genesis, nil
}
