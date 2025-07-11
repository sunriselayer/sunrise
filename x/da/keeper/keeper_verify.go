package keeper

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

// This file contains the logic for the data verification and lifecycle management.
// It was previously handled in the EndBlocker.

func (k Keeper) DeleteRejectedDataOvertime(ctx sdk.Context, duration time.Duration) error {
	rejectedData, err := k.GetSpecificStatusDataBeforeTime(ctx, types.Status_STATUS_REJECTED, ctx.BlockTime().Add(-duration).Unix())
	if err != nil {
		return err
	}
	for _, data := range rejectedData {
		if data.Status == types.Status_STATUS_REJECTED {
			err = k.DeletePublishedData(ctx, data)
			if err != nil {
				return err
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
				return err
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
				return err
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
			thresholdDec, err := math.LegacyNewDecFromStr(threshold) // TODO: remove with Dec
			if err != nil {
				return err
			}
			invalidityThreshold := thresholdDec.MulInt64(int64(len(data.ShardDoubleHashes)))
			if math.LegacyNewDec(int64(len(invalidIndices))).GTE(invalidityThreshold) {
				// Get active validators
				activeValidators := []string{}
				iterator, err := k.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
				if err != nil {
					return err
				}
				defer iterator.Close()
				for ; iterator.Valid(); iterator.Next() {
					validator, err := k.StakingKeeper.Validator(ctx, iterator.Value())
					if err != nil {
						return err
					}
					if validator.IsBonded() {
						activeValidators = append(activeValidators, validator.GetOperator())
					}
				}
				data.ChallengingValidators = activeValidators

				data.Status = types.Status_STATUS_CHALLENGING
				data.Timestamp = ctx.BlockTime()
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (k Keeper) ChangeToVerifiedFromChallengePeriod(ctx sdk.Context, duration time.Duration) error {
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
				return err
			}
			// refunds collateral to the publisher
			publisher := sdk.MustAccAddressFromBech32(data.Publisher)
			err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, publisher, data.PublishDataCollateral)
			if err != nil {
				return err
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

	replicationFactorDec, err := math.LegacyNewDecFromStr(replicationFactor) // TODO: remove with Dec
	if err != nil {
		return err
	}
	faultValidators := make(map[string]sdk.ValAddress)

	for _, data := range challengingData {
		if data.Status == types.Status_STATUS_CHALLENGING {
			proofs, err := k.GetProofs(ctx, data.MetadataUri)
			if err != nil {
				return err
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
				return err
			}
			indexedValidators := make(map[int64][]sdk.ValAddress)
			for _, valAddrStr := range data.ChallengingValidators {
				valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
				if err != nil {
					return err
				}
				indices := types.ShardIndicesForValidator(valAddr, int64(threshold), int64(len(data.ShardDoubleHashes)))
				for _, index := range indices {
					indexedValidators[index] = append(indexedValidators[index], valAddr)
				}
			}

			safeShardIndices := []int64{}
			for index := 0; int64(index) < int64(len(data.ShardDoubleHashes)); index++ {
				i := int64(index)
				if len(data.ShardDoubleHashes) < int(data.ParityShardCount) {
					return types.ErrInvalidShardCounts
				}

				// replication_factor_with_parity = replication_factor * data_shard_count / (data_shard_count + parity_shard_count)
				replicationFactorWithParity := replicationFactorDec.
					MulInt64(int64(len(data.ShardDoubleHashes) - int(data.ParityShardCount))).
					QuoInt64(int64(len(data.ShardDoubleHashes)))

				// threshold_of_proofs_per_shard = ceil(replication_factor_with_parity * 2 / 3)
				thresholdOfProofsPerShard := replicationFactorWithParity.
					MulInt64(2).
					QuoInt64(3).
					Ceil().
					TruncateInt64()

				numAssignedValidators := int64(len(indexedValidators[i]))
				proofCount := shardProofCount[i]

				var requiredProofs int64
				if numAssignedValidators < thresholdOfProofsPerShard {
					// Use 2/3 of assigned validators as the threshold if it's smaller
					requiredProofs = math.LegacyNewDec(numAssignedValidators).
						MulInt64(2).
						QuoInt64(3).
						Ceil().
						TruncateInt64()
				} else {
					requiredProofs = thresholdOfProofsPerShard
				}

				// A shard is safe if the number of proofs is greater than or equal to the required number.
				if proofCount >= requiredProofs {
					safeShardIndices = append(safeShardIndices, i)
					if numAssignedValidators > 0 {
						for _, valAddr := range indexedValidators[i] {
							if !shardProofSubmitted[i][sdk.AccAddress(valAddr).String()] {
								faultValidators[valAddr.String()] = valAddr
							}
						}
					}
				}
			}

			// valid_shards < data_shard_count
			invalidities, err := k.GetInvalidities(ctx, data.MetadataUri)
			if err != nil {
				return err
			}
			if int64(len(safeShardIndices))+int64(data.ParityShardCount) < int64(len(data.ShardDoubleHashes)) {
				data.Status = types.Status_STATUS_REJECTED
				data.Timestamp = ctx.BlockTime()
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					return err
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
						return err
					}
				}
			} else {
				data.Status = types.Status_STATUS_VERIFIED
				data.Timestamp = ctx.BlockTime()
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					return err
				}

				publisherRefund := data.PublishDataCollateral
				for _, invalidity := range invalidities {
					challenger := sdk.MustAccAddressFromBech32(invalidity.Sender)
					if checkCorrectInvalidity(invalidity, safeShardIndices) {
						// if all shards in the invalidity are missing, refund to the challenger
						err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, challenger, data.SubmitInvalidityCollateral)
						if err != nil {
							return err
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
					return err
				}
			}

			// Count challenge & fault validators
			err = k.SetChallengeCounter(ctx, k.GetChallengeCounter(ctx)+1)
			if err != nil {
				return err
			}
			for _, valAddr := range faultValidators {
				count, err := k.GetFaultCounter(ctx, valAddr)
				if err != nil {
					return err
				}
				err = k.SetFaultCounter(ctx, valAddr, count+1)
				if err != nil {
					return err
				}
			}

			// Clean up proofs data
			for _, proof := range proofs {
				addr, err := k.addressCodec.StringToBytes(proof.Sender)
				if err != nil {
					return err
				}
				err = k.DeleteProof(ctx, proof.MetadataUri, addr)
				if err != nil {
					return err
				}
			}

			// Clean up invalidity data
			for _, invalidity := range invalidities {
				addr, err := k.addressCodec.StringToBytes(invalidity.Sender)
				if err != nil {
					return err
				}
				err = k.DeleteInvalidity(ctx, invalidity.MetadataUri, addr)
				if err != nil {
					return err
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
