package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k Keeper) CreateEpoch(ctx sdk.Context, previousEpochId, epochId uint64) error {
	gauges, err := k.Tally(ctx)
	if err != nil {
		return err
	}

	if len(gauges) == 0 {
		return nil
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
