package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k Keeper) CreateEpoch(ctx sdk.Context, epochId uint64) error {
	// Finalize bribe for new epoch & remove old epochs
	if err := k.FinalizeBribeForEpoch(ctx, epochId); err != nil {
		return err
	}

	// Tally voting power to create gauges and delete votes
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
		StartTime:  ctx.BlockTime().Unix(),
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
