package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
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
				// Calculate relative weight
				relativeWeight := weight.Quo(poolTotalWeights[poolWeight.PoolId])
				// Save BribeAllocation
				bribeAllocation := types.BribeAllocation{
					Address:         voterStr,
					EpochId:         epochId,
					PoolId:          poolWeight.PoolId,
					Weight:          relativeWeight.String(),
					ClaimedBribeIds: []uint64{},
				}
				err = k.SetBribeAllocation(ctx, bribeAllocation)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// FinalizeBribeForEpoch saves vote weights for bribes distribution and processes unclaimed bribes for old epochs that have passed their claim period
func (k Keeper) FinalizeBribeForEpoch(ctx sdk.Context, epochId uint64) error {
	// Save vote weights
	if err := k.SaveVoteWeightsForBribes(ctx, epochId); err != nil {
		return err
	}

	// Process unclaimed bribes for old epochs that have passed their claim period
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// Use EpochBlocks as the claim period for now
	expiredEpochId := k.GetBribeExpiredEpochId(ctx)
	if epochId > params.BribeClaimEpochs {
		newExpiredEpochId := epochId - params.BribeClaimEpochs
		if newExpiredEpochId > expiredEpochId {
			for id := expiredEpochId; id < newExpiredEpochId; id++ {
				if err := k.ProcessUnclaimedBribes(ctx, id); err != nil {
					return err
				}
				// Remove epoch
				if err := k.RemoveEpoch(ctx, id); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// ProcessUnclaimedBribes processes unclaimed bribes for an epoch that has passed its claim period
func (k Keeper) ProcessUnclaimedBribes(ctx context.Context, epochId uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Process all bribes for this epoch
	totalUnclaimed := sdk.NewCoins()
	var bribesToRemove []uint64

	err := k.Bribes.Indexes.EpochId.Walk(
		ctx,
		collections.NewPrefixedPairRange[uint64, uint64](epochId),
		func(_ uint64, bribeId uint64) (bool, error) {
			bribe, found, err := k.GetBribe(ctx, bribeId)
			if !found {
				return true, nil
			}
			if err != nil {
				return false, err
			}
			// Calculate unclaimed amount
			unclaimedAmount := bribe.Amount.Sub(bribe.ClaimedAmount...)
			if !unclaimedAmount.IsZero() {
				totalUnclaimed = totalUnclaimed.Add(unclaimedAmount...)
			}
			bribesToRemove = append(bribesToRemove, bribeId)
			return false, nil
		})
	if err != nil {
		return err
	}

	// If there are no bribes to remove, return nil
	if len(bribesToRemove) == 0 {
		return nil
	}

	// Remove all bribes for this epoch
	for _, bribeId := range bribesToRemove {
		if err := k.RemoveBribe(ctx, bribeId); err != nil {
			continue
		}
	}

	// Use the new index to get keys to remove
	iter, err := k.BribeAllocations.Indexes.EpochId.MatchExact(ctx, epochId)
	if err != nil {
		return err
	}
	defer iter.Close()

	var keysToRemove []collections.Triple[sdk.AccAddress, uint64, uint64]
	for ; iter.Valid(); iter.Next() {
		key, err := iter.PrimaryKey()
		if err != nil {
			return err
		}
		keysToRemove = append(keysToRemove, key)
	}

	for _, key := range keysToRemove {
		if err := k.RemoveBribeAllocation(ctx, key); err != nil {
			continue
		}
	}

	// If there are unclaimed bribes, send them to fee collector
	if !totalUnclaimed.IsZero() {
		feeCollectorAddr := k.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			sdkCtx,
			types.BribeAccount,
			feeCollectorAddr,
			totalUnclaimed,
		); err != nil {
			return errorsmod.Wrap(err, "failed to send unclaimed bribes to fee collector")
		}
	}

	// Set the expired epoch id
	if err := k.SetBribeExpiredEpochId(ctx, epochId); err != nil {
		return err
	}

	return nil
}
