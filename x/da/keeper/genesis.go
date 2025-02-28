package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, data := range genState.PublishedData {
		if err := k.SetPublishedData(ctx, data); err != nil {
			return err
		}
	}
	for _, proof := range genState.Proofs {
		if err := k.SetProof(ctx, proof); err != nil {
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

	genesis.PublishedData = k.GetAllPublishedData(ctx)
	genesis.Proofs, err = k.GetAllProofs(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
