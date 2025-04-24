package keeper

import (
	"context"
	"fmt"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, acc := range genState.LockupAccounts {
		err := k.SetLockupAccount(ctx, acc)
		if err != nil {
			return fmt.Errorf("%w: %s #%d", err, acc.Owner, acc.Id)
		}
	}

	for index, msgInit := range genState.InitLockupMsgs {
		err := k.InitLockupAccountFromMsg(ctx, msgInit)
		if err != nil {
			return fmt.Errorf("invalid genesis account msg init at index %d, msg %s: %w", index, msgInit, err)
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

	genesis.LockupAccounts, err = k.GetAllLockupAccounts(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
