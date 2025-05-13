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
	// Set all the votes
	for _, vote := range genState.Votes {
		err := k.SetVote(ctx, vote)
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

	genesis.Epochs, err = k.GetAllEpoch(ctx)
	if err != nil {
		return nil, err
	}
	genesis.EpochCount, err = k.GetEpochCount(ctx)
	if err != nil {
		return nil, err
	}
	genesis.Votes, err = k.GetAllVotes(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
