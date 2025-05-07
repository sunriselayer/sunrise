package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) EndBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := k.DeleteExpiredBlobDeclarations(sdkCtx)
	if err != nil {
		k.Logger().Error("failed to delete expired blob declarations", "error", err)
	}

	err = k.DeleteExpiredBlobCommitments(sdkCtx)
	if err != nil {
		k.Logger().Error("failed to delete expired blob commitments", "error", err)
	}

	err = k.DeleteUnusedValidatorsPowerSnapshots(sdkCtx)
	if err != nil {
		k.Logger().Error("failed to delete unused validators power snapshots", "error", err)
	}

	return nil
}

func (k Keeper) DeleteExpiredBlobDeclarations(ctx sdk.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	now := sdkCtx.BlockTime()

	err := k.BlobDeclarations.Indexes.Expiry.Walk(
		ctx,
		collections.NewPrefixUntilPairRange[int64, []byte](now.Unix()),
		func(_ int64, shardsMerkleRoot []byte) (stop bool, err error) {
			err = k.BlobDeclarations.Remove(ctx, shardsMerkleRoot)
			if err != nil {
				return false, err
			}

			return false, err
		},
	)

	return err
}

func (k Keeper) DeleteExpiredBlobCommitments(ctx sdk.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	now := sdkCtx.BlockTime()

	err := k.BlobCommitments.Indexes.Expiry.Walk(
		ctx,
		collections.NewPrefixUntilPairRange[int64, []byte](now.Unix()),
		func(_ int64, shardsMerkleRoot []byte) (stop bool, err error) {
			err = k.BlobCommitments.Remove(ctx, shardsMerkleRoot)
			if err != nil {
				return false, err
			}

			// Remove related challenges
			challenges, err := k.GetAllChallengesByShardsMerkleRoot(ctx, shardsMerkleRoot)
			if err != nil {
				return false, err
			}

			for _, challenge := range challenges {
				err = k.Challenges.Remove(ctx, challenge.Id)
				if err != nil {
					return false, err
				}
			}

			return false, nil
		},
	)

	return err
}

func (k Keeper) CountBlobDeclarationsWithHeight(ctx sdk.Context, blockHeight int64) (uint64, error) {
	count := uint64(0)
	err := k.BlobDeclarations.Indexes.BlockHeight.Walk(
		ctx,
		collections.NewPrefixedPairRange[int64, []byte](blockHeight),
		func(_ int64, declarationKey []byte) (stop bool, err error) {
			count++
			return false, nil
		},
	)

	return count, err
}

func (k Keeper) CheckAddValidatorsPowerSnapshot(ctx sdk.Context) error {
	count, err := k.CountBlobDeclarationsWithHeight(ctx, ctx.BlockHeight())
	if err != nil {
		return err
	}

	if count > 0 {
		return k.TakeValidatorsPowerSnapshot(ctx)
	}

	return nil
}

func (k Keeper) DeleteUnusedValidatorsPowerSnapshots(ctx sdk.Context) error {
	err := k.ValidatorsPowerSnapshots.Walk(
		ctx,
		nil,
		func(blockHeight int64, _ types.ValidatorsPowerSnapshot) (stop bool, err error) {
			count, err := k.CountBlobDeclarationsWithHeight(ctx, blockHeight)
			if err != nil {
				return false, err
			}

			if count > 0 {
				return true, nil
			}

			err = k.ValidatorsPowerSnapshots.Remove(ctx, blockHeight)
			if err != nil {
				return false, err
			}

			return false, nil
		},
	)

	return err
}
