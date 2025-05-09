package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k Keeper) CreateEpoch(ctx sdk.Context, previousEpochId, epochId uint64) error {
	results, err := k.Tally(ctx)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		return nil
	}

	// Create new gauges
	gauges := []types.Gauge{}
	for _, result := range results {
		gauge := types.Gauge{
			PreviousEpochId: previousEpochId,
			PoolId:          result.PoolId,
			Count:           result.Count,
		}
		err := k.SetGauge(ctx, gauge)
		if err != nil {
			return err
		}
		gauges = append(gauges, gauge)
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	epoch := types.Epoch{
		Id:         epochId,
		StartBlock: ctx.BlockHeight(),
		EndBlock:   ctx.BlockHeight() + params.EpochBlocks,
		Gauges:     gauges,
	}
	err = k.SetEpoch(ctx, epoch)
	if err != nil {
		return err
	}

	// Emit event
	if err := ctx.EventManager().EmitTypedEvent(&types.EventStartNewEpoch{
		EpochId:    epoch.Id,
		StartBlock: epoch.StartBlock,
		EndBlock:   epoch.EndBlock,
	}); err != nil {
		return err
	}

	return nil
}
