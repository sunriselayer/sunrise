package keeper

import (
	"context"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) EndBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params, err := k.Params.Get(ctx)
	if err != nil {
		k.Logger().Error("failed to get params", "error", err)
		return err
	}

	// If STATUS_REJECTED is overtime, remove from the store
	err = k.DeleteRejectedDataOvertime(sdkCtx, params.RejectedRemovalPeriod)
	if err != nil {
		k.Logger().Error("failed to delete rejected data overtime", "error", err)
	}

	// IF STATUS_VERIFIED is overtime, remove from store
	err = k.DeleteVerifiedDataOvertime(sdkCtx, params.VerifiedRemovalPeriod)
	if err != nil {
		k.Logger().Error("failed to delete verified data overtime", "error", err)
	}

	// if STATUS_CHALLENGE_PERIOD receives invalidity above the threshold, change to STATUS_CHALLENGING
	err = k.ChangeToChallengingFromChallengePeriod(sdkCtx, params.ChallengeThreshold)
	if err != nil {
		k.Logger().Error("failed to change to challenging from challenge period", "error", err)
	}

	// if STATUS_CHALLENGE_PERIOD is expired, change to STATUS_VERIFIED
	err = k.ChangeToVerifiedFromProofPeriod(sdkCtx, params.ChallengePeriod)
	if err != nil {
		k.Logger().Error("failed to change to verified from proof period", "error", err)
	}

	// if STATUS_CHALLENGING, tally validity_proofs
	err = k.TallyValidityProofs(sdkCtx, params.ProofPeriod, params.ReplicationFactor)
	if err != nil {
		k.Logger().Error("failed to tally validity proofs", "error", err)
	}

	// slash epoch moved from PreBlocker
	if sdkCtx.BlockHeight()%int64(params.SlashEpoch) == 0 {
		k.HandleSlashEpoch(sdkCtx)
	}

	return nil
}

func (k Keeper) DeleteRejectedDataOvertime(ctx sdk.Context, duration time.Duration) error {
	rejectedData, err := k.GetSpecificStatusDataBeforeTime(ctx, types.Status_STATUS_REJECTED, ctx.BlockTime().Add(-duration).Unix())
	if err != nil {
		return err
	}
	for _, data := range rejectedData {
		if data.Status == types.Status_STATUS_REJECTED {
			err = k.DeletePublishedData(ctx, data)
			if err != nil {
				k.Logger().Error("failed to delete published data", "metadata_uri", data.MetadataUri, "error", err)
				continue
			}
		}
	}
	return nil
}

func (k Keeper) DeleteVerifiedDataOvertime(ctx sdk.Context, duration time.Duration) error {
	verifiedData, err := k.GetSpecificStatusDataBeforeTime(ctx, types.Status_STATUS_VERIFIED, ctx.BlockTime().Add(-duration).Unix())
	if err != nil {
		return err
	}
	for _, data := range verifiedData {
		if data.Status == types.Status_STATUS_VERIFIED {
			err = k.DeletePublishedData(ctx, data)
			if err != nil {
				k.Logger().Error("failed to delete published data", "metadata_uri", data.MetadataUri, "error", err)
				continue
			}
		}
	}
	return nil
}

func (k Keeper) ChangeToChallengingFromChallengePeriod(ctx sdk.Context, threshold string) error {
	challengePeriodData, err := k.GetSpecificStatusData(ctx, types.Status_STATUS_CHALLENGE_PERIOD)
	if err != nil {
		return err
	}
	for _, data := range challengePeriodData {
		if data.Status == types.Status_STATUS_CHALLENGE_PERIOD {
			invalidities, err := k.GetInvalidities(ctx, data.MetadataUri)
			if err != nil {
				k.Logger().Error("failed to get invalidity", "metadata_uri", data.MetadataUri, "error", err)
				continue
			}
			seen := make(map[int64]bool)
			invalidIndices := []int64{}
			for _, invalidity := range invalidities {
				for _, index := range invalidity.Indices {
					if _, ok := seen[index]; !ok {
						seen[index] = true
						invalidIndices = append(invalidIndices, index)
					}
				}
			}
			thresholdDec := math.LegacyMustNewDecFromStr(threshold) // TODO: remove with Dec
			invalidityThreshold := thresholdDec.MulInt64(int64(len(data.ShardDoubleHashes)))
			if math.LegacyNewDec(int64(len(invalidIndices))).GTE(invalidityThreshold) {
				data.Status = types.Status_STATUS_CHALLENGING
				data.Timestamp = ctx.BlockTime()
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					k.Logger().Error("failed to set published data", "metadata_uri", data.MetadataUri, "error", err)
					continue
				}
			}
		}
	}
	return nil
}

func (k Keeper) ChangeToVerifiedFromProofPeriod(ctx sdk.Context, duration time.Duration) error {
	expiredChallengePeriodData, err := k.GetSpecificStatusDataBeforeTime(ctx, types.Status_STATUS_CHALLENGE_PERIOD, ctx.BlockTime().Add(-duration).Unix())
	if err != nil {
		return err
	}
	for _, data := range expiredChallengePeriodData {
		if data.Status == types.Status_STATUS_CHALLENGE_PERIOD {
			data.Status = types.Status_STATUS_VERIFIED
			data.Timestamp = ctx.BlockTime()
			err = k.SetPublishedData(ctx, data)
			if err != nil {
				k.Logger().Error("failed to set published data", "metadata_uri", data.MetadataUri, "error", err)
				continue
			}
			// refunds collateral to the publisher
			publisher := sdk.MustAccAddressFromBech32(data.Publisher)
			err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, publisher, data.PublishDataCollateral)
			if err != nil {
				k.Logger().Error("failed to send coins to publisher", "publisher", publisher.String(), "metadata_uri", data.MetadataUri, "error", err)
				continue
			}
		}
	}
	return nil
}

func (k Keeper) TallyValidityProofs(ctx sdk.Context, duration time.Duration, replicationFactor string) error {
	challengingData, err := k.GetSpecificStatusDataBeforeTime(ctx, types.Status_STATUS_CHALLENGING, ctx.BlockTime().Add(-duration).Unix())
	if err != nil {
		return err
	}

	activeValidators := []sdk.ValAddress{}
	iterator, err := k.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		return err
	}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validator, err := k.StakingKeeper.Validator(ctx, iterator.Value())
		if err != nil {
			k.Logger().Error("failed to get validator", "error", err)
			continue
		}
		if validator.IsBonded() {
			activeValidators = append(activeValidators, sdk.ValAddress(iterator.Value()))
		}
	}

	replicationFactorDec := math.LegacyMustNewDecFromStr(replicationFactor) // TODO: remove with Dec
	faultValidators := make(map[string]sdk.ValAddress)

	for _, data := range challengingData {
		if data.Status == types.Status_STATUS_CHALLENGING {
			proofs, err := k.GetProofs(ctx, data.MetadataUri)
			if err != nil {
				k.Logger().Error("failed to get proofs", "metadata_uri", data.MetadataUri, "error", err)
				continue
			}
			shardProofCount := make(map[int64]int64)
			shardProofSubmitted := make(map[int64]map[string]bool)
			for _, proof := range proofs {
				for _, index := range proof.Indices {
					shardProofCount[index]++
					if shardProofSubmitted[index] == nil {
						shardProofSubmitted[index] = make(map[string]bool)
					}
					shardProofSubmitted[index][proof.Sender] = true
				}
			}

			threshold, err := k.GetZkpThreshold(ctx, uint64(len(data.ShardDoubleHashes)))
			if err != nil {
				k.Logger().Error("failed to get zkp threshold", "metadata_uri", data.MetadataUri, "error", err)
				continue
			}
			indexedValidators := make(map[int64][]sdk.ValAddress)
			for _, valAddr := range activeValidators {
				indices := types.ShardIndicesForValidator(valAddr, int64(threshold), int64(len(data.ShardDoubleHashes)))
				for _, index := range indices {
					indexedValidators[index] = append(indexedValidators[index], valAddr)
				}
			}

			safeShardIndices := []int64{}
			for index, proofCount := range shardProofCount {
				if len(data.ShardDoubleHashes) < int(data.ParityShardCount) {
					k.Logger().Error("parity shard count is greater than total shard count", "metadata_uri", data.MetadataUri)
					continue
				}
				// replication_factor_with_parity = replication_factor * data_shard_count / (data_shard_count + parity_shard_count)
				replicationFactorWithParity := replicationFactorDec.
					MulInt64(int64(len(data.ShardDoubleHashes) - int(data.ParityShardCount))).
					QuoInt64(int64(len(data.ShardDoubleHashes)))

				// len(zkp_including_this_shard) / replication_factor_with_parity >= 2/3
				if math.LegacyNewDec(proofCount).GTE(
					replicationFactorWithParity.
						MulInt64(2).
						QuoInt64(3)) {
					safeShardIndices = append(safeShardIndices, index)
					for _, valAddr := range indexedValidators[index] {
						if !shardProofSubmitted[index][sdk.AccAddress(valAddr).String()] {
							faultValidators[valAddr.String()] = valAddr
						}
					}
				}
			}

			// valid_shards < data_shard_count
			invalidities, err := k.GetInvalidities(ctx, data.MetadataUri)
			if err != nil {
				k.Logger().Error("failed to get invalidity", "metadata_uri", data.MetadataUri, "error", err)
				continue
			}
			if int64(len(safeShardIndices))+int64(data.ParityShardCount) < int64(len(data.ShardDoubleHashes)) {
				data.Status = types.Status_STATUS_REJECTED
				data.Timestamp = ctx.BlockTime()
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					k.Logger().Error("failed to set published data", "metadata_uri", data.MetadataUri, "error", err)
					continue
				}

				// distribute publish collateral to challengers as a reward.
				publishCollateral := data.PublishDataCollateral
				reward := sdk.Coins{}
				for _, coin := range publishCollateral {
					dividedAmount := math.LegacyNewDecFromInt(coin.Amount).QuoInt64(int64(len(invalidities))).TruncateInt()
					reward = append(reward, sdk.NewCoin(coin.Denom, dividedAmount))
				}

				// rewards collateral + reward to challengers
				for _, invalidity := range invalidities {
					challenger := sdk.MustAccAddressFromBech32(invalidity.Sender)
					err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, challenger, data.SubmitInvalidityCollateral.Add(reward...))
					if err != nil {
						k.Logger().Error("failed to send coins to challenger", "challenger", challenger.String(), "metadata_uri", data.MetadataUri, "error", err)
						continue
					}
				}
			} else {
				data.Status = types.Status_STATUS_VERIFIED
				data.Timestamp = ctx.BlockTime()
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					k.Logger().Error("failed to set published data", "metadata_uri", data.MetadataUri, "error", err)
					continue
				}

				publisherRefund := data.PublishDataCollateral
				for _, invalidity := range invalidities {
					challenger := sdk.MustAccAddressFromBech32(invalidity.Sender)
					if checkCorrectInvalidity(invalidity, safeShardIndices) {
						// if all shards in the invalidity are missing, refund to the challenger
						err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, challenger, data.SubmitInvalidityCollateral)
						if err != nil {
							k.Logger().Error("failed to send coins to challenger", "challenger", challenger.String(), "metadata_uri", data.MetadataUri, "error", err)
							continue
						}
					} else {
						// if at least one safe shard is included, pass to the publisher
						publisherRefund = publisherRefund.Add(data.SubmitInvalidityCollateral...)
					}
				}

				// refunds publish collateral + fault challengers' collateral to the publisher
				publisher := sdk.MustAccAddressFromBech32(data.Publisher)
				err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, publisher, publisherRefund)
				if err != nil {
					k.Logger().Error("failed to send coins to publisher", "publisher", publisher.String(), "metadata_uri", data.MetadataUri, "error", err)
					continue
				}
			}

			// Count challenge & fault validators
			err = k.SetChallengeCounter(ctx, k.GetChallengeCounter(ctx)+1)
			if err != nil {
				k.Logger().Error("failed to set challenge counter", "error", err)
				continue
			}
			for _, valAddr := range faultValidators {
				count, err := k.GetFaultCounter(ctx, valAddr)
				if err != nil {
					k.Logger().Error("failed to get fault counter", "validator", valAddr.String(), "error", err)
					continue
				}
				err = k.SetFaultCounter(ctx, valAddr, count+1)
				if err != nil {
					k.Logger().Error("failed to set fault counter", "validator", valAddr.String(), "error", err)
					continue
				}
			}

			// Clean up proofs data
			for _, proof := range proofs {
				addr, err := k.addressCodec.StringToBytes(proof.Sender)
				if err != nil {
					k.Logger().Error("failed to convert sender to bytes", "sender", proof.Sender, "error", err)
					continue
				}
				err = k.DeleteProof(ctx, proof.MetadataUri, addr)
				if err != nil {
					k.Logger().Error("failed to delete proof", "metadata_uri", proof.MetadataUri, "sender", proof.Sender, "error", err)
					continue
				}
			}

			// Clean up invalidity data
			for _, invalidity := range invalidities {
				addr, err := k.addressCodec.StringToBytes(invalidity.Sender)
				if err != nil {
					k.Logger().Error("failed to convert sender to bytes", "sender", invalidity.Sender, "error", err)
					continue
				}
				err = k.DeleteInvalidity(ctx, invalidity.MetadataUri, addr)
				if err != nil {
					k.Logger().Error("failed to delete invalidity", "metadata_uri", invalidity.MetadataUri, "sender", invalidity.Sender, "error", err)
					continue
				}
			}
		}
	}
	return nil
}

func checkCorrectInvalidity(invalidity types.Invalidity, safeShardIndices []int64) bool {
	safeIndexMap := make(map[int64]bool)
	for _, index := range safeShardIndices {
		safeIndexMap[index] = true
	}
	for _, index := range invalidity.Indices {
		if _, ok := safeIndexMap[index]; ok {
			return false // includes a safe shard
		}
	}
	return true // all of them is invalid
}
