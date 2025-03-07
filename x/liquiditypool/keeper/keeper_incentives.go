package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AllocateIncentive(ctx sdk.Context, poolId uint64, sender sdk.AccAddress, incentiveCoins sdk.Coins) error {
	pool, err := k.getPoolForSwap(ctx, poolId)
	if err != nil {
		return err
	}

	feeAccumulator, err := k.GetFeeAccumulator(ctx, poolId)
	if err != nil {
		return err
	}
	liquidity, err := math.LegacyNewDecFromStr(pool.CurrentTickLiquidity)
	if err != nil {
		return err
	}
	feeGrowth := sdk.NewDecCoinsFromCoins(incentiveCoins...).QuoDecTruncate(liquidity)
	err = k.AddToAccumulator(ctx, feeAccumulator, feeGrowth)
	if err != nil {
		return err
	}

	return k.bankKeeper.SendCoins(ctx, sender, pool.GetFeesAddress(), incentiveCoins)
}
