package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
	return nil
}

func (k Keeper) BeginBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	// Transfer a portion of inflation rewards from fee collector to `x/liquidityincentive` pool.
	feeCollector := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	incentiveModule := authtypes.NewModuleAddress(types.ModuleName)
	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return err
	}
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return err
	}
	feeBalance := k.bankKeeper.GetBalance(ctx, feeCollector, feeDenom)
	amount := feeBalance.Amount

	k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, types.ModuleName, sdk.NewCoins(feeBalance))

	// Convert fee denom to bond denom
	err = k.tokenConverterKeeper.ConvertReverse(ctx, amount, incentiveModule)
	if err != nil {
		return err
	}

	// Distribute incentives to gauges
	lastEpoch, found, err := k.GetLastEpoch(ctx)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	amountDec := math.LegacyNewDecFromInt(amount)
	totalCount := math.LegacyZeroDec()
	for _, gauge := range lastEpoch.Gauges {
		totalCount = totalCount.Add(math.LegacyNewDecFromInt(gauge.Count))
	}

	if totalCount.IsZero() {
		return nil
	}
	for _, gauge := range lastEpoch.Gauges {
		weight := math.LegacyNewDecFromInt(gauge.Count).Quo(totalCount)
		allocationDec := amountDec.Mul(weight)
		if allocationDec.IsPositive() {
			err := k.liquidityPoolKeeper.AllocateIncentive(
				ctx,
				gauge.PoolId,
				incentiveModule,
				sdk.NewCoins(sdk.NewCoin(bondDenom, allocationDec.TruncateInt())),
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
	lastEpoch, found, err := k.GetLastEpoch(ctx)
	if err != nil {
		return err
	}
	if !found {
		err := k.CreateEpoch(ctx, 0, 1)
		if err != nil {
			ctx.Logger().Error("epoch creation error", err)
			return nil
		}
	} else if ctx.BlockHeight() >= lastEpoch.EndBlock {
		// End current epoch and start new one
		if err := k.FinalizeBribeForEpoch(ctx); err != nil {
			ctx.Logger().Error("epoch ending error", err)
			return nil
		}

		err := k.CreateEpoch(ctx, lastEpoch.Id, lastEpoch.Id+1)
		if err != nil {
			ctx.Logger().Error("epoch creation error", err)
			return nil
		}
		// remove old epoch and gauges
		epochs, err := k.GetAllEpoch(ctx)
		if err != nil {
			return err
		}
		if len(epochs) > 2 {
			epoch := epochs[0]
			err := k.RemoveEpoch(ctx, epoch.Id)
			if err != nil {
				return err
			}
			for _, gauge := range epoch.Gauges {
				err := k.RemoveGauge(ctx, gauge.PreviousEpochId, gauge.PoolId)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
