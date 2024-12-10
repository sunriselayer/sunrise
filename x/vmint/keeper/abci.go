package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/sunriselayer/sunrise/x/vmint/types"
)

var (
	millisecondsPerYear = int64(31556952000)
)

// Refer https://github.com/cosmos/cosmos-sdk/blob/v0.50.10/x/mint/abci.go
func (k Keeper) BeginBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	params := k.GetParams(ctx)

	supplyBond := k.bankKeeper.GetSupply(ctx, params.BondDenom)
	supplyFee := k.bankKeeper.GetSupply(ctx, params.FeeDenom)
	totalSupply := supplyBond.Amount.Add(supplyFee.Amount)

	inflationRateCapInitial, err := math.LegacyNewDecFromStr(params.InflationRateCapInitial)
	if err != nil {
		return err
	}
	inflationRateCapMinimum, err := math.LegacyNewDecFromStr(params.InflationRateCapMinimum)
	if err != nil {
		return err
	}
	disinflationRate, err := math.LegacyNewDecFromStr(params.DisinflationRate)
	if err != nil {
		return err
	}

	annualProvision := types.CalculateAnnualProvision(
		ctx,
		inflationRateCapInitial,
		inflationRateCapMinimum,
		disinflationRate,
		params.SupplyCap,
		params.Genesis,
		totalSupply,
	)

	lastBlockTime, err := k.LastBlockTime.Get(ctx)
	if err != nil {
		lastBlockTime = params.Genesis
	}

	blockTime := ctx.BlockTime()
	duration := blockTime.Sub(lastBlockTime)

	// annualProvision * (durationMilliseconds / millisecondsPerYear)
	blockProvision := annualProvision.Mul(math.NewInt(duration.Milliseconds())).Quo(math.NewInt(millisecondsPerYear))

	if blockProvision.IsPositive() {
		coins := sdk.NewCoins(sdk.NewCoin(params.BondDenom, blockProvision))
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, coins)
		if err != nil {
			return err
		}
	}

	// set
	err = k.LastBlockTime.Set(ctx, blockTime)
	if err != nil {
		return err
	}

	return nil
}
