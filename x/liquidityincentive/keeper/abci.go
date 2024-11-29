package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sunriselayer/sunrise/app/consts"
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
	epoch := types.Epoch{
		Id:         epochId,
		StartBlock: ctx.BlockHeight(),
		EndBlock:   ctx.BlockHeight() + params.EpochBlocks,
		Gauges:     gauges,
	}
	k.SetEpoch(ctx, epoch)
	return nil
}

func (k Keeper) BeginBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	// Transfer a portion of inflation rewards from fee collector to `x/liquidityincentive` pool.
	feeCollector := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	fees := k.bankKeeper.GetAllBalances(ctx, feeCollector)
	vRiseAmount := fees.AmountOf(consts.BondDenom)
	amount := sdk.NewCoin(consts.BondDenom, vRiseAmount)
	feesDec := sdk.NewDecCoinsFromCoins(amount)

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
		ratio := weight.Ratio.Quo(totalWeight)
		allocationDec := incentiveFeesDec.MulDecTruncate(ratio)
		allocation, _ := allocationDec.TruncateDecimal()
		if allocation.IsAllPositive() {
			err := k.liquidityPoolKeeper.AllocateIncentive(
				ctx,
				weight.PoolId,
				authtypes.NewModuleAddress(authtypes.FeeCollectorName),
				allocation,
			)
			if err != nil {
				ctx.Logger().Error("failure in incentive allocation", "error", err)
			}
		}
	}

	return nil
}

func (k Keeper) EndBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	// Create a new `Epoch` if the last `Epoch` has ended or the first `Epoch` has not been created.
	lastEpoch, found := k.GetLastEpoch(ctx)
	if !found {
		err := k.CreateEpoch(ctx, 0, 1)
		if err != nil {
			ctx.Logger().Error("epoch creation error", err)
			return nil
		}
	} else if ctx.BlockHeight() >= lastEpoch.EndBlock {
		err := k.CreateEpoch(ctx, lastEpoch.Id, lastEpoch.Id+1)
		if err != nil {
			ctx.Logger().Error("epoch creation error", err)
			return nil
		}
		// remove old epoch and gauges
		epochs := k.GetAllEpoch(ctx)
		if len(epochs) > 2 {
			epoch := epochs[0]
			k.RemoveEpoch(ctx, epoch.Id)
			for _, gauge := range epoch.Gauges {
				k.RemoveGauge(ctx, gauge.PreviousEpochId, gauge.PoolId)
			}
		}
	}
	return nil
}
