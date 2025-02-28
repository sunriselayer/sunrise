package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetChallengeCounter returns the challenge counter
func (k Keeper) GetChallengeCounter(ctx context.Context) (uint64, error) {
	val, err := k.ChallengeCounts.Get(ctx)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// SetChallengeCounter sets the challenge counter
func (k Keeper) SetChallengeCounter(ctx context.Context, count uint64) error {
	return k.ChallengeCounts.Set(ctx, count)
}

// GetFaultCounter returns the fault counter for a validator
func (k Keeper) GetFaultCounter(ctx context.Context, operator sdk.ValAddress) (uint64, error) {
	has, err := k.FaultCounts.Has(ctx, operator)
	if err != nil {
		return 0, err
	}

	if !has {
		return 0, nil
	}

	return k.FaultCounts.Get(ctx, operator)
}

// SetFaultCounter sets the fault counter for a validator
func (k Keeper) SetFaultCounter(ctx context.Context, operator sdk.ValAddress, faultCounter uint64) error {
	return k.FaultCounts.Set(ctx, operator, faultCounter)
}

// DeleteFaultCounter removes the fault counter for a validator
func (k Keeper) DeleteFaultCounter(ctx context.Context, operator sdk.ValAddress) error {
	return k.FaultCounts.Remove(ctx, operator)
}

// IterateFaultCounters iterates over all fault counters
func (k Keeper) IterateFaultCounters(ctx context.Context,
	handler func(operator sdk.ValAddress, faultCount uint64) (stop bool, err error),
) error {
	return k.FaultCounts.Walk(
		ctx,
		nil,
		func(key []byte, value uint64) (bool, error) {
			stop, err := handler(key, value)
			if err != nil {
				return false, err
			}
			return stop, nil
		},
	)
}

// HandleSlashEpoch handles the slash epoch
func (k Keeper) HandleSlashEpoch(ctx sdk.Context) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		k.Logger.Error("failed to get params", "error", err)
		return
	}
	slashFaultThreshold := math.LegacyMustNewDecFromStr(params.SlashFaultThreshold)
	slashFraction := math.LegacyMustNewDecFromStr(params.SlashFraction)

	challengeCount, err := k.GetChallengeCounter(ctx)
	if err != nil {
		k.Logger.Error("failed to get challenge counter", "error", err)
		return
	}

	// reset counter
	if err := k.SetChallengeCounter(ctx, 0); err != nil {
		k.Logger.Error("failed to reset challenge counter", "error", err)
		return
	}

	threshold := slashFaultThreshold.MulInt64(int64(challengeCount)).TruncateInt().Uint64()
	powerReduction := k.StakingKeeper.PowerReduction(ctx)

	err = k.IterateFaultCounters(ctx, func(operator sdk.ValAddress, faultCount uint64) (bool, error) {
		validator, err := k.StakingKeeper.Validator(ctx, operator)
		if err != nil {
			k.Logger.Error("failed to get validator", "error", err, "operator", operator.String())
			return false, nil
		}

		if err := k.DeleteFaultCounter(ctx, operator); err != nil {
			k.Logger.Error("failed to delete fault counter", "error", err, "operator", operator.String())
			return false, nil
		}

		if validator.IsJailed() || !validator.IsBonded() {
			return false, nil
		}

		if faultCount <= threshold {
			return false, nil
		}

		consAddr, err := validator.GetConsAddr()
		if err != nil {
			k.Logger.Error("failed to get consensus address", "error", err, "operator", operator.String())
			return false, nil
		}

		if err := k.SlashingKeeper.Slash(
			ctx, consAddr, slashFraction,
			validator.GetConsensusPower(powerReduction),
			ctx.BlockHeight()-sdk.ValidatorUpdateDelay-1,
		); err != nil {
			k.Logger.Error("failed to slash validator", "error", err, "operator", operator.String())
			return false, nil
		}

		if err := k.SlashingKeeper.Jail(ctx, consAddr); err != nil {
			k.Logger.Error("failed to jail validator", "error", err, "operator", operator.String())
			return false, nil
		}

		return false, nil
	})
	if err != nil {
		k.Logger.Error("failed to iterate fault counters", "error", err)
	}
}
