package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Set params
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	// Set markets
	for _, market := range genState.Markets {
		if err := k.Markets.Set(ctx, market.Denom, market); err != nil {
			return err
		}
	}

	// Set user positions
	for _, position := range genState.UserPositions {
		if err := k.UserPositions.Set(ctx, collections.Join(position.UserAddress, position.Denom), position); err != nil {
			return err
		}
	}

	// Set borrows
	for _, borrow := range genState.Borrows {
		if err := k.Borrows.Set(ctx, borrow.Id, borrow); err != nil {
			return err
		}
	}

	// Set borrow id sequence
	if err := k.BorrowId.Set(ctx, genState.BorrowCount); err != nil {
		return err
	}

	return nil
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	genesis := types.DefaultGenesis()

	// Get params
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	genesis.Params = params

	// Export markets
	err = k.Markets.Walk(ctx, nil, func(denom string, market types.Market) (stop bool, err error) {
		genesis.Markets = append(genesis.Markets, market)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	// Export user positions
	err = k.UserPositions.Walk(ctx, nil, func(key collections.Pair[string, string], position types.UserPosition) (stop bool, err error) {
		genesis.UserPositions = append(genesis.UserPositions, position)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	// Export borrows
	err = k.Borrows.Walk(ctx, nil, func(id uint64, borrow types.Borrow) (stop bool, err error) {
		genesis.Borrows = append(genesis.Borrows, borrow)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	// Get borrow count
	borrowCount, err := k.BorrowId.Peek(ctx)
	if err != nil {
		// If no borrows have been created yet, the sequence might not exist
		genesis.BorrowCount = 0
	} else {
		genesis.BorrowCount = borrowCount
	}

	return genesis, nil
}
