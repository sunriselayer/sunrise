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
				// Process only pools with bribes
				bribe, found, err := k.GetBribeByEpochAndPool(ctx, epochId, poolWeight.PoolId)
				if err != nil {
					return err
				}
				if !found {
					continue
				}

				// Calculate relative weight
				relativeWeight := weight.Quo(poolTotalWeights[poolWeight.PoolId])

				// Save UnclaimedBribe
				unclaimedBribe := types.UnclaimedBribe{
					Address: voterStr,
					BribeId: bribe.Id,
					Weight:  relativeWeight.String(),
				}

				voterAddr, err := k.addressCodec.StringToBytes(voterStr)
				if err != nil {
					return err
				}
				key := collections.Join(sdk.AccAddress(voterAddr), bribe.Id)
				if err := k.UnclaimedBribes.Set(ctx, key, unclaimedBribe); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// FinalizeBribeForEpoch ends the current epoch and starts a new one
func (k Keeper) FinalizeBribeForEpoch(ctx sdk.Context) error {

	// Get the ID of the epoch to end
	currentEpochId, err := k.GetEpochCount(ctx)
	if err != nil {
		return err
	}

	// Save vote weights
	if err := k.SaveVoteWeightsForBribes(ctx, currentEpochId); err != nil {
		ctx.Logger().Error("failed to save vote weights for bribes",
			"epoch_id", currentEpochId,
			"error", err,
		)
	}

	// Process unclaimed bribes for old epochs that have passed their claim period
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// Use EpochBlocks as the claim period for now
	if currentEpochId > uint64(params.EpochBlocks) {
		epochToProcess := currentEpochId - uint64(params.EpochBlocks)
		if err := k.ProcessUnclaimedBribes(ctx, epochToProcess); err != nil {
			// Only log the error and continue with epoch ending process
			ctx.Logger().Error("failed to process unclaimed bribes",
				"epoch_id", epochToProcess,
				"error", err,
			)
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
			bribe, _, err := k.GetBribe(ctx, bribeId)
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

	// Remove all bribes for this epoch
	for _, bribeId := range bribesToRemove {
		if err := k.RemoveBribe(ctx, bribeId); err != nil {
			continue
		}
	}

	// Remove all unclaimed bribes for this epoch
	var keysToRemove []collections.Pair[sdk.AccAddress, uint64]
	err = k.UnclaimedBribes.Walk(ctx,
		nil, // Iterate over all unclaimed bribes with Pair key
		func(key collections.Pair[sdk.AccAddress, uint64], value types.UnclaimedBribe) (bool, error) {
			if value.EpochId == epochId { // Check bribe's epochId
				keysToRemove = append(keysToRemove, key)
			}
			return false, nil
		},
	)
	if err != nil {
		return err
	}

	// Remove all unclaimed bribes for this epoch
	for _, key := range keysToRemove {
		if err := k.UnclaimedBribes.Remove(ctx, key); err != nil {
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

	return nil
}
