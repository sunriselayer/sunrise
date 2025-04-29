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

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	cacheCtx, write := sdk.UnwrapSDKContext(ctx).CacheContext()

	// Transfer a portion of inflation rewards from fee collector to `x/liquidityincentive` pool.
	feeCollector := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	incentiveModule := authtypes.NewModuleAddress(types.ModuleName)
	feeDenom, err := k.feeKeeper.FeeDenom(cacheCtx)
	if err != nil {
		k.Logger().Error("failed to get fee denom", "error", err)
		return
	}
	bondDenom, err := k.stakingKeeper.BondDenom(cacheCtx)
	if err != nil {
		k.Logger().Error("failed to get bond denom", "error", err)
		return
	}

	// Check the Gauge count is not zero.
	// Distribute incentives to gauges
	lastEpoch, found, err := k.GetLastEpoch(cacheCtx)
	if err != nil {
		k.Logger().Error("failed to get last epoch", "error", err)
		return
	}
	if !found {
		return
	}

	totalCount := math.LegacyZeroDec()
	for _, gauge := range lastEpoch.Gauges {
		totalCount = totalCount.Add(math.LegacyNewDecFromInt(gauge.Count))
	}

	if totalCount.IsZero() {
		k.Logger().Info("total count is zero")
		return
	}

	// Send a portion of inflation rewards from fee collector to `x/liquidityincentive` pool.
	feeBalance := k.bankKeeper.GetBalance(cacheCtx, feeCollector, feeDenom)
	feeCollectorAmountDec := math.LegacyNewDecFromInt(feeBalance.Amount)
	params, err := k.Params.Get(cacheCtx)
	if err != nil {
		k.Logger().Error("failed to get params", "error", err)
		return
	}
	stakingRewardRatioDec, err := math.LegacyNewDecFromStr(params.StakingRewardRatio)
	if err != nil {
		k.Logger().Error("failed to parse staking reward ratio", "error", err)
		return
	}
	incentiveAmount := feeCollectorAmountDec.Mul(math.LegacyOneDec().Sub(stakingRewardRatioDec)).TruncateInt()
	err = k.bankKeeper.SendCoinsFromModuleToModule(cacheCtx, authtypes.FeeCollectorName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(feeDenom, incentiveAmount)))
	if err != nil {
		k.Logger().Error("failed to send coins from fee collector to liquidity incentive module", "error", err)
		return
	}

	// Convert fee denom to bond denom in the `x/liquidityincentive` module account.
	err = k.tokenConverterKeeper.ConvertReverse(cacheCtx, incentiveAmount, incentiveModule)
	if err != nil {
		k.Logger().Error("failed to convert fee denom to bond denom", "error", err)
		return
	}

	// Get `x/liquidityincentive` module's incentive balance.
	incentiveBalance := k.bankKeeper.GetBalance(cacheCtx, incentiveModule, bondDenom)

	// Distribute incentives to gauges
	for _, gauge := range lastEpoch.Gauges {
		weight := math.LegacyNewDecFromInt(gauge.Count).Quo(totalCount)
		allocationDec := math.LegacyNewDecFromInt(incentiveBalance.Amount).Mul(weight)
		if allocationDec.IsPositive() {
			err := k.liquidityPoolKeeper.AllocateIncentive(
				cacheCtx,
				gauge.PoolId,
				incentiveModule,
				sdk.NewCoins(sdk.NewCoin(bondDenom, allocationDec.TruncateInt())),
			)
			if err != nil {
				k.Logger().Error("failure in incentive allocation", "error", err)
			}
		}
	}

	write()
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
