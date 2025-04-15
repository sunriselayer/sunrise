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
	votes, err := k.GetAllVotes(ctx)
	if err != nil {
		return err
	}

	for _, vote := range votes {
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
	}

	// Save relative weights for each voter
	for _, vote := range votes {
		voterStr := vote.Sender

		for _, poolWeight := range vote.PoolWeights {
			weight, err := math.LegacyNewDecFromStr(poolWeight.Weight)
			if err != nil {
				continue
			}

			if weight.IsPositive() && !poolTotalWeights[poolWeight.PoolId].IsZero() {
				// Process only pools with bribes
				bribeKey := collections.Join(epochId, poolWeight.PoolId)
				bribe, err := k.Bribes.Get(ctx, bribeKey)
				if err != nil {
					return err
				}
				if bribe == (types.Bribe{}) {
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

				key := collections.Join3(sdk.MustAccAddressFromBech32(voterStr), epochId, poolWeight.PoolId)
				if err := k.UnclaimedBribes.Set(ctx, key, unclaimedBribe); err != nil {
					return err
				}
			}
		}
	}

	return nil
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
	epoch, err := k.Epochs.Get(ctx, epochId)
	if err != nil {
		return err
	}
	if epoch.Id == 0 {
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
	err = k.UnclaimedBribes.Walk(ctx, collections.NewPrefixedTripleRange[sdk.AccAddress, uint64, uint64](sdk.AccAddress{}),
		func(key collections.Triple[sdk.AccAddress, uint64, uint64], unclaimed types.UnclaimedBribe) (bool, error) {
			if key.K2() == epochId {
				return false, k.UnclaimedBribes.Remove(ctx, key)
			}
			return false, nil
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
	}

	return nil
}
