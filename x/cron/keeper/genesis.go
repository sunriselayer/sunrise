package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/cron/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Set all the schedules
	for _, elem := range genState.ScheduleList {
		err := k.AddSchedule(sdkCtx, elem.Name, elem.Period, elem.Msgs, elem.ExecutionStage)
		if err != nil {
			panic(err)
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
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	genesis.ScheduleList, err = k.GetAllSchedules(sdkCtx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
