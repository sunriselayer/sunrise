package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) GetChallengeCounter(ctx context.Context) uint64 {
	val, err := k.ChallengeCounts.Get(ctx)
	if err != nil {
		return 0
	}

	return val
}

func (k Keeper) SetChallengeCounter(ctx context.Context, count uint64) error {
	err := k.ChallengeCounts.Set(ctx, count)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetFaultCounter(ctx context.Context, operator sdk.ValAddress) (count uint64, err error) {
	has, err := k.FaultCounts.Has(ctx, operator)
	if err != nil {
		return 0, err
	}

	if !has {
		return 0, nil
	}

	count, err = k.FaultCounts.Get(ctx, operator)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (k Keeper) SetFaultCounter(ctx context.Context, operator sdk.ValAddress, faultCounter uint64) error {
	err := k.FaultCounts.Set(ctx, operator, faultCounter)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteFaultCounter(ctx context.Context, operator sdk.ValAddress) error {
	err := k.FaultCounts.Remove(ctx, operator)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) IterateFaultCounters(ctx context.Context,
	handler func(operator sdk.ValAddress, faultCount uint64) (stop bool, err error),
) error {
	return k.FaultCounts.Walk(
		ctx,
		nil,
		func(key []byte, value uint64) (bool, error) {
			return handler(key, value)
		},
	)
}

func (k Keeper) HandleSlashEpoch(ctx sdk.Context) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	slashFaultThreshold, err := math.LegacyNewDecFromStr(params.SlashFaultThreshold) // TODO: remove with Dec
	if err != nil {
		return err
	}
	slashFraction, err := math.LegacyNewDecFromStr(params.SlashFraction) // TODO: remove with Dec
	if err != nil {
		return err
	}
	challengeCount := k.GetChallengeCounter(ctx)
	// reset counter
	err = k.SetChallengeCounter(ctx, 0)
	if err != nil {
		return err
	}
	threshold := slashFaultThreshold.MulInt64(int64(challengeCount)).Ceil().TruncateInt().Uint64()
	powerReduction := k.StakingKeeper.PowerReduction(ctx)
	return k.IterateFaultCounters(ctx, func(operator sdk.ValAddress, faultCount uint64) (bool, error) {
		validator, err := k.StakingKeeper.Validator(ctx, operator)
		if err != nil {
			if errors.Is(err, stakingtypes.ErrNoValidatorFound) {
				k.Logger().Info("validator not found during slash epoch, skipping", "validator", operator.String())
				return false, nil
			}
			return false, err
		}

		defer func() {
			err := k.DeleteFaultCounter(ctx, operator)
			if err != nil {
				k.Logger().Error("failed to delete fault counter", "error", err, "validator", operator.String())
			}
		}()
		if validator.IsJailed() || !validator.IsBonded() {
			return false, nil
		}

		if faultCount <= threshold {
			return false, nil
		}

		consAddr, err := validator.GetConsAddr()
		if err != nil {
			return false, err
		}

		err = k.SlashingKeeper.Slash(
			ctx, consAddr, slashFraction,
			validator.GetConsensusPower(powerReduction),
			ctx.BlockHeight()-sdk.ValidatorUpdateDelay-1,
		)
		if err != nil {
			return false, err
		}
		err = k.SlashingKeeper.Jail(ctx, consAddr)
		if err != nil {
			return false, err
		}
		return false, nil
	})
}
