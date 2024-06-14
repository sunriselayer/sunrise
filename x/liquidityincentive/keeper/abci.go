package keeper

import (
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k Keeper) CreateEpoch(ctx sdk.Context, previousEpochId, epochId uint64) error {
	weights, err := k.Tally(ctx)
	if err != nil {
		return err
	}

	if len(weights) == 0 {
		return nil
	}

	gauges := []types.Gauge{}
	for _, weight := range weights {
		gauge := types.Gauge{
			PreviousEpochId: previousEpochId,
			PoolId:          weight.PoolId,
			Ratio:           weight.Weight,
		}
		k.SetGauge(ctx, gauge)
		gauges = append(gauges, gauge)
	}

	params := k.GetParams(ctx)
	k.SetEpoch(ctx, types.Epoch{
		Id:         epochId,
		StartBlock: ctx.BlockHeight(),
		EndBlock:   ctx.BlockHeight() + params.EpochBlocks,
		Gauges:     gauges,
	})
	return nil
}

// BeginBlocker sets the proposer for determining distribution during endblock
// and distribute rewards for the previous block.
func (k Keeper) BeginBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// Create a new `Epoch` if the last `Epoch` has ended or the first `Epoch` has not been created.
	lastEpoch, found := k.GetLastEpoch(ctx)
	if !found {
		k.CreateEpoch(ctx, 0, 1)
	} else if lastEpoch.EndBlock >= ctx.BlockHeight() {
		k.CreateEpoch(ctx, lastEpoch.Id, lastEpoch.Id+1)
		// TODO: remove old epoch
	}
	return nil
}

// EndBlocker called every block, process inflation, update validator set.
func (k Keeper) EndBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Transfer a portion of inflation rewards from fee collector to `x/liquidityincentive` pool.
	feeCollector := k.authKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	fees := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesDec := sdk.NewDecCoinsFromCoins(fees...)

	params := k.GetParams(ctx)
	incentiveFeesDec := feesDec.MulDecTruncate(math.LegacyOneDec().Sub(params.StakingRewardRatio))

	lastEpoch, found := k.GetLastEpoch(ctx)
	if !found {
		return nil
	}

	totalWeight := math.LegacyZeroDec()
	for _, weight := range lastEpoch.Gauges {
		totalWeight = totalWeight.Add(weight.Ratio)
	}

	if totalWeight.IsZero() {
		return nil
	}
	for _, weight := range lastEpoch.Gauges {
		allocationDec := incentiveFeesDec.MulDecTruncate(weight.Ratio.Quo(totalWeight))
		allocation, _ := allocationDec.TruncateDecimal()

		err := k.liquidityPoolKeeper.AllocateIncentive(
			ctx,
			weight.PoolId,
			authtypes.NewModuleAddress(authtypes.FeeCollectorName),
			allocation,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
