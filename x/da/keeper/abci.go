package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) EndBlocker(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params, err := k.Params.Get(ctx)
	if err != nil {
		k.Logger.Error(err.Error())
		return
	}
	// if STATUS_CHALLENGE_PERIOD receives invalidity above the threshold, change to STATUS_CHALLENGING
	challengePeriodData, err := k.GetSpecificStatusDataBeforeTime(sdkCtx, types.Status_STATUS_CHALLENGE_PERIOD, sdkCtx.BlockTime().Unix())
	if err != nil {
		k.Logger.Error(err.Error())
		return
	}
	for _, data := range challengePeriodData {
		if data.Status == types.Status_STATUS_CHALLENGE_PERIOD {
			invalidities, err := k.GetInvalidities(sdkCtx, data.MetadataUri)
			if err != nil {
				k.Logger.Error("failed to get invalidities", "error", err)
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
			threshold := math.LegacyMustNewDecFromStr(params.ChallengeThreshold).MulInt64(int64(len(data.ShardDoubleHashes)))
			if math.LegacyNewDec(int64(len(invalidIndices))).GTE(threshold) {
				data.Status = types.Status_STATUS_CHALLENGING
				data.ChallengeTimestamp = sdkCtx.BlockTime()
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					k.Logger.Error(err.Error())
					return
				}
			}
		}
	}

	// if STATUS_CHALLENGE_PERIOD is expired, change to STATUS_VERIFIED
	expiredChallengePeriodData, err := k.GetSpecificStatusDataBeforeTime(sdkCtx, types.Status_STATUS_CHALLENGE_PERIOD, sdkCtx.BlockTime().Add(-params.ChallengePeriod).Unix())
	if err != nil {
		k.Logger.Error(err.Error())
		return
	}
	for _, data := range expiredChallengePeriodData {
		if data.Status == types.Status_STATUS_CHALLENGE_PERIOD {
			data.Status = types.Status_STATUS_VERIFIED
			err = k.SetPublishedData(ctx, data)
			if err != nil {
				k.Logger.Error(err.Error())
				return
			}
			// refunds collateral to the publisher
			publisher := sdk.MustAccAddressFromBech32(data.Publisher)
			err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, publisher, data.PublishDataCollateral)
			if err != nil {
				k.Logger.Error(err.Error())
				return
			}
		}
	}

	// if STATUS_CHALLENGING, tally validity_proofs
	challengingData, err := k.GetSpecificStatusDataBeforeTime(sdkCtx, types.Status_STATUS_CHALLENGING, sdkCtx.BlockTime().Add(-params.ChallengePeriod-params.ProofPeriod).Unix())
	if err != nil {
		k.Logger.Error(err.Error())
		return
	}

	activeValidators := []sdk.ValAddress{}
	iterator, err := k.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		k.Logger.Error(err.Error())
		return
	}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validator, err := k.StakingKeeper.Validator(ctx, iterator.Value())
		if err != nil {
			k.Logger.Error(err.Error())
			return
		}
		if validator.IsBonded() {
			activeValidators = append(activeValidators, sdk.ValAddress(iterator.Value()))
		}
	}

	replicationFactor := math.LegacyMustNewDecFromStr(params.ReplicationFactor) // TODO: remove with Dec
	faultValidators := make(map[string]sdk.ValAddress)

	for _, data := range challengingData {
		if data.Status == types.Status_STATUS_CHALLENGING {
			proofs, err := k.GetProofs(sdkCtx, data.MetadataUri)
			if err != nil {
				k.Logger.Error(err.Error())
				continue
			}
			shardProofCount := make(map[int64]int64)
			shardProofSubmitted := make(map[int64]map[string]bool)
			for _, proof := range proofs {
				for _, index := range proof.Indices {
					shardProofCount[index]++
					shardProofSubmitted[index][proof.Sender] = true
				}
			}

			threshold, err := k.GetZkpThreshold(ctx, uint64(len(data.ShardDoubleHashes)))
			if err != nil {
				k.Logger.Error(err.Error())
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
					k.Logger.Error("parity shard count is greater than total shard count")
					continue
				}
				// replication_factor_with_parity = replication_factor * data_shard_count / (data_shard_count + parity_shard_count)
				replicationFactorWithParity := replicationFactor.
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
			invalidities, err := k.GetInvalidities(sdkCtx, data.MetadataUri)
			if err != nil {
				k.Logger.Error("failed to get invalidities", "error", err)
				continue
			}
			if int64(len(safeShardIndices))+int64(data.ParityShardCount) < int64(len(data.ShardDoubleHashes)) {
				data.Status = types.Status_STATUS_REJECTED
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					k.Logger.Error(err.Error())
					return
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
						k.Logger.Error(err.Error())
						continue
					}
				}
			} else {
				data.Status = types.Status_STATUS_VERIFIED
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					k.Logger.Error(err.Error())
					return
				}

				publisherRefund := data.PublishDataCollateral
				for _, invalidity := range invalidities {
					challenger := sdk.MustAccAddressFromBech32(invalidity.Sender)
					if checkCorrectInvalidity(invalidity, safeShardIndices) {
						// if all shards in the invalidity are missing, refund to the challenger
						err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, challenger, data.SubmitInvalidityCollateral)
						if err != nil {
							k.Logger.Error(err.Error())
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
					k.Logger.Error(err.Error())
					return
				}
			}

			// Count challenge & fault validators
			err = k.SetChallengeCounter(ctx, k.GetChallengeCounter(ctx)+1)
			if err != nil {
				k.Logger.Error(err.Error())
				return
			}
			for _, valAddr := range faultValidators {
				count, err := k.GetFaultCounter(ctx, valAddr)
				if err != nil {
					k.Logger.Error(err.Error())
					continue
				}
				err = k.SetFaultCounter(ctx, valAddr, count+1)
				if err != nil {
					k.Logger.Error(err.Error())
					continue
				}
			}

			// Clean up proofs data
			for _, proof := range proofs {
				addr, err := k.addressCodec.StringToBytes(proof.Sender)
				if err != nil {
					k.Logger.Error(err.Error())
					continue
				}
				err = k.DeleteProof(sdkCtx, proof.MetadataUri, addr)
				if err != nil {
					k.Logger.Error(err.Error())
					continue
				}
			}

			// Clean up invalidity data
			for _, invalidity := range invalidities {
				addr, err := k.addressCodec.StringToBytes(invalidity.Sender)
				if err != nil {
					k.Logger.Error(err.Error())
					continue
				}
				err = k.DeleteInvalidity(sdkCtx, invalidity.MetadataUri, addr)
				if err != nil {
					k.Logger.Error(err.Error())
					continue
				}
			}
		}
	}

	// If STATUS_REJECTED is overtime, remove from the store
	rejectedData, err := k.GetSpecificStatusDataBeforeTime(sdkCtx, types.Status_STATUS_REJECTED, sdkCtx.BlockTime().Add(-params.RejectedRemovalPeriod).Unix())
	if err != nil {
		k.Logger.Error(err.Error())
		return
	}
	for _, data := range rejectedData {
		if data.Status == types.Status_STATUS_REJECTED {
			err = k.DeletePublishedData(sdkCtx, data)
			if err != nil {
				k.Logger.Error(err.Error())
				continue
			}
		}
	}

	// If VerifiedRemovalPeriod id positive and STATUS_VERIFIED is overtime, remove from store
	if params.VerifiedRemovalPeriod > 0 {
		verifiedData, err := k.GetSpecificStatusDataBeforeTime(sdkCtx, types.Status_STATUS_VERIFIED, sdkCtx.BlockTime().Add(-params.VerifiedRemovalPeriod).Unix())
		if err != nil {
			k.Logger.Error(err.Error())
			return
		}
		for _, data := range verifiedData {
			if data.Status == types.Status_STATUS_VERIFIED {
				err = k.DeletePublishedData(sdkCtx, data)
				if err != nil {
					k.Logger.Error(err.Error())
					continue
				}
			}
		}
	}

	// slash epoch moved from vote_extension
	if sdkCtx.BlockHeight()%int64(params.SlashEpoch) == 0 {
		k.HandleSlashEpoch(sdkCtx)
	}
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
