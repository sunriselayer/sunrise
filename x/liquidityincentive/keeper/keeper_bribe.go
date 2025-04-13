package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// SaveVoteWeightsForBribes saves vote weights for bribes distribution
func (k Keeper) SaveVoteWeightsForBribes(ctx context.Context, epochId uint64) error {
	// Calculate total vote weights for each pool
	poolTotalWeights := make(map[uint64]math.LegacyDec)

	// Process all votes
	err := k.Votes.Walk(ctx, collections.NewPrefixedRange[sdk.AccAddress](), func(voter sdk.AccAddress, vote types.Vote) (bool, error) {
		for _, poolWeight := range vote.PoolWeights {
			weight, err := math.LegacyNewDecFromStr(poolWeight.Weight)
			if err != nil {
				continue
			}

			if weight.IsPositive() {
				if _, ok := poolTotalWeights[poolWeight.PoolId]; !ok {
					poolTotalWeights[poolWeight.PoolId] = math.LegacyZeroDec()
				}

				poolTotalWeights[poolWeight.PoolId] = poolTotalWeights[poolWeight.PoolId].Add(weight)
			}
		}
		return false, nil
	})
	if err != nil {
		return err
	}

	// Save relative weights for each voter
	err = k.Votes.Walk(ctx, collections.NewPrefixedRange[sdk.AccAddress](), func(voter sdk.AccAddress, vote types.Vote) (bool, error) {
		voterStr, err := k.addressCodec.BytesToString(voter)
		if err != nil {
			return false, err
		}

		for _, poolWeight := range vote.PoolWeights {
			weight, err := math.LegacyNewDecFromStr(poolWeight.Weight)
			if err != nil {
				continue
			}

			if weight.IsPositive() && !poolTotalWeights[poolWeight.PoolId].IsZero() {
				// Process only pools with bribes
				bribeKey := collections.Join(epochId, poolWeight.PoolId)
				_, found, err := k.Bribes.Get(ctx, bribeKey)
				if err != nil {
					return false, err
				}
				if !found {
					continue
				}

				// Calculate relative weight
				relativeWeight := weight.Quo(poolTotalWeights[poolWeight.PoolId])

				// Save UnclaimedBribe
				unclaimedBribe := types.UnclaimedBribe{
					Address: voterStr,
					EpochId: epochId,
					PoolId:  poolWeight.PoolId,
					Weight:  relativeWeight.String(),
				}

				key := collections.Join3(voterStr, epochId, poolWeight.PoolId)
				if err := k.UnclaimedBribes.Set(ctx, key, unclaimedBribe); err != nil {
					return false, err
				}
			}
		}
		return false, nil
	})

	return err
}

// EndEpoch ends the current epoch and starts a new one
func (k Keeper) EndEpoch(ctx context.Context) error {
	// ... existing code for ending the current epoch ...

	// Get the ID of the epoch to end
	currentEpochId := k.GetCurrentEpochId(ctx)

	// Save vote weights
	if err := k.SaveVoteWeightsForBribes(ctx, currentEpochId); err != nil {
		k.Logger(sdk.UnwrapSDKContext(ctx)).Error(
			"failed to save vote weights for bribes",
			"epoch_id", currentEpochId,
			"error", err,
		)
	}

	// Process unclaimed bribes for old epochs that have passed their claim period
	if currentEpochId > k.GetParams(ctx).BribeClaimEpochs {
		epochToProcess := currentEpochId - k.GetParams(ctx).BribeClaimEpochs
		if err := k.ProcessUnclaimedBribes(ctx, epochToProcess); err != nil {
			// Only log the error and continue with epoch ending process
			k.Logger(sdk.UnwrapSDKContext(ctx)).Error(
				"failed to process unclaimed bribes",
				"epoch_id", epochToProcess,
				"error", err,
			)
		}
	}

	// ... rest of the epoch ending code ...

	return nil
}

// ProcessUnclaimedBribes processes unclaimed bribes for an epoch that has passed its claim period
func (k Keeper) ProcessUnclaimedBribes(ctx context.Context, epochId uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check if the epoch exists
	_, found, err := k.Epochs.Get(ctx, epochId)
	if err != nil {
		return err
	}
	if !found {
		return errorsmod.Wrapf(types.ErrEpochNotFound, "epoch %d not found", epochId)
	}

	// Process all bribes for this epoch
	totalUnclaimed := sdk.NewCoins()

	err = k.Bribes.Walk(ctx, collections.NewPrefixedPairRange[uint64, uint64](epochId),
		func(key collections.Pair[uint64, uint64], bribe types.Bribe) (bool, error) {
			// Calculate unclaimed amount
			unclaimedAmount := bribe.Amount.Sub(bribe.ClaimedAmount)
			if !unclaimedAmount.IsZero() {
				totalUnclaimed = totalUnclaimed.Add(unclaimedAmount)
			}

			// Remove all unclaimed bribes for this epoch
			if err := k.Bribes.Remove(ctx, key); err != nil {
				return false, err
			}

			return false, nil
		})

	if err != nil {
		return err
	}

	// Remove all unclaimed bribes for this epoch
	err = k.UnclaimedBribes.Walk(ctx, collections.NewPrefixedTripleRange[string, uint64, uint64]("", epochId, 0),
		func(key collections.Triple[string, uint64, uint64], unclaimed types.UnclaimedBribe) (bool, error) {
			return false, k.UnclaimedBribes.Remove(ctx, key)
		})

	if err != nil {
		return err
	}

	// If there are unclaimed bribes, send them to fee collector
	if !totalUnclaimed.IsZero() {
		feeCollectorAddr := k.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			sdkCtx,
			types.ModuleName,
			feeCollectorAddr,
			totalUnclaimed,
		); err != nil {
			return errorsmod.Wrap(err, "failed to send unclaimed bribes to fee collector")
		}

		// Emit event
		if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventUnclaimedBribesProcessed{
			EpochId: epochId,
			Amount:  totalUnclaimed,
		}); err != nil {
			return err
		}
	}

	return nil
}
