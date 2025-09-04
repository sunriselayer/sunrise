package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenfactory/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if genState.Params.DenomCreationFee == nil {
		genState.Params.DenomCreationFee = sdk.NewCoins()
	}
	err := k.Params.Set(ctx, genState.Params)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, genDenom := range genState.GetFactoryDenoms() {
		creator, _, err := types.DeconstructDenom(genDenom.GetDenom())
		if err != nil {
			return err
		}
		err = k.createDenomAfterValidation(sdkCtx, creator, genDenom.GetDenom())
		if err != nil {
			return err
		}
		err = k.setAuthorityMetadata(sdkCtx, genDenom.GetDenom(), genDenom.GetAuthorityMetadata())
		if err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = k.DenomFromCreator.Walk(ctx, nil, func(key collections.Pair[sdk.AccAddress, string], value []byte) (stop bool, err error) {
		denom := key.K2()
		authorityMetadata, err := k.GetAuthorityMetadata(sdkCtx, denom)
		if err != nil {
			return true, err
		}

		genesis.FactoryDenoms = append(genesis.FactoryDenoms, types.GenesisDenom{
			Denom:             denom,
			AuthorityMetadata: authorityMetadata,
		})
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
