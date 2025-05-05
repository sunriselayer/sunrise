package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, deputy := range genState.Deputies {
		validator, err := k.StakingKeeper.ValidatorAddressCodec().StringToBytes(deputy.Validator)
		if err != nil {
			return err
		}
		if err := k.SetDeputy(ctx, validator, deputy.Address); err != nil {
			return err
		}
	}
	for _, declaration := range genState.BlobDeclarations {
		if err := k.SetBlobDeclaration(ctx, declaration); err != nil {
			return err
		}
	}
	for _, snapshot := range genState.ValidatorPowerSnapshots {
		if err := k.SetValidatorPowerSnapshot(ctx, snapshot); err != nil {
			return err
		}
	}
	for _, commitment := range genState.BlobCommitments {
		if err := k.SetBlobCommitment(ctx, commitment); err != nil {
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

	genesis.Deputies, err = k.GetAllDeputies(ctx)
	if err != nil {
		return nil, err
	}
	genesis.BlobDeclarations, err = k.GetAllBlobDeclarations(ctx)
	if err != nil {
		return nil, err
	}
	genesis.ValidatorPowerSnapshots, err = k.GetAllValidatorPowerSnapshots(ctx)
	if err != nil {
		return nil, err
	}
	genesis.BlobCommitments, err = k.GetAllBlobCommitments(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
