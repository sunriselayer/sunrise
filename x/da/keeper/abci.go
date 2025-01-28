package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) EndBlocker(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// TODO: error handling
	params, _ := k.Params.Get(ctx)
	replicationFactor := math.LegacyMustNewDecFromStr(params.ReplicationFactor) // TODO: remove with Dec
	challengePeriodData, err := k.GetUnverifiedDataBeforeTime(sdkCtx, sdkCtx.BlockTime().Add(-params.ChallengePeriod).Unix())
	if err != nil {
		k.Logger.Error(err.Error())
		return
	}

	for _, data := range challengePeriodData {
		if data.Status == types.Status_STATUS_VOTE_EXTENSION {
			data.Status = types.Status_STATUS_VERIFIED
		}
		if err = k.SetPublishedData(ctx, data); err != nil {
			return
		}

		publisher := sdk.MustAccAddressFromBech32(data.Publisher)
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, publisher, data.Collateral)
		if err != nil {
			k.Logger.Error(err.Error())
			return
		}
	}

	proofPeriodData, err := k.GetUnverifiedDataBeforeTime(sdkCtx, sdkCtx.BlockTime().Add(-params.ChallengePeriod-params.ProofPeriod).Unix())
	if err != nil {
		k.Logger.Error(err.Error())
		return
	}

	activeValidators := []sdk.ValAddress{}
	numActiveValidators := int64(0)
	// votingPowers := make(map[string]int64)
	// powerReduction := k.StakingKeeper.PowerReduction(ctx)
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
			// valAddrStr := validator.GetOperator()
			// valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			// if err != nil {
			// 	k.Logger().Error(err.Error())
			// 	return
			// }

			// votingPowers[sdk.AccAddress(valAddr).String()] = validator.GetConsensusPower(powerReduction)
			numActiveValidators++
		}
	}

	faultValidators := make(map[string]sdk.ValAddress)

	for _, data := range proofPeriodData {
		if data.Status == types.Status_STATUS_CHALLENGE_FOR_FRAUD {
			// bondedTokens, err := k.StakingKeeper.TotalBondedTokens(ctx)
			// if err != nil {
			// 	k.Logger().Error(err.Error())
			// 	return
			// }

			// totalBondedPower := sdk.TokensToConsensusPower(bondedTokens, k.StakingKeeper.PowerReduction(ctx))
			// thresholdPower := params.VoteThreshold.MulInt64(totalBondedPower).RoundInt().Int64()
			proofs := k.GetProofs(sdkCtx, data.MetadataUri)
			shardProofCount := make(map[int64]int64)
			shardProofSubmitted := make(map[int64]map[string]bool)
			for _, proof := range proofs {
				for _, index := range proof.Indices {
					shardProofCount[index]++
					shardProofSubmitted[index][proof.Sender] = true
				}
			}

			threshold := k.GetZkpThreshold(ctx, uint64(len(data.ShardDoubleHashes)))
			indexedValidators := make(map[int64][]sdk.ValAddress)
			for _, valAddr := range activeValidators {
				indices := types.ShardIndicesForValidator(valAddr, int64(threshold), int64(len(data.ShardDoubleHashes)))
				for _, index := range indices {
					indexedValidators[index] = append(indexedValidators[index], valAddr)
				}
			}

			safeShardCount := int64(0)
			for index, proofCount := range shardProofCount {
				// replication_factor_with_parity = replication_factor * data_shard_count / (data_shard_count + parity_shard_count)
				replicationFactorWithParity := replicationFactor.
					MulInt64(int64(len(data.ShardDoubleHashes) - int(data.ParityShardCount))).
					QuoInt64(int64(len(data.ShardDoubleHashes)))

				// len(zkp_including_this_shard) / replication_factor_with_parity >= 2/3
				if math.LegacyNewDec(proofCount).GTE(
					replicationFactorWithParity.
						MulInt64(2).
						QuoInt64(3)) {
					safeShardCount++
					for _, valAddr := range indexedValidators[index] {
						if !shardProofSubmitted[index][sdk.AccAddress(valAddr).String()] {
							faultValidators[valAddr.String()] = valAddr
						}
					}
				}
			}

			// valid_shards < data_shard_count
			if safeShardCount+int64(data.ParityShardCount) < int64(len(data.ShardDoubleHashes)) {
				// TODO: might require rejected records as well
				data.Status = types.Status_STATUS_REJECTED
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					k.Logger.Error(err.Error())
					return
				}
				// k.DeletePublishedData(sdkCtx, data)

				challenger := sdk.MustAccAddressFromBech32(data.Challenger)
				err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, challenger, data.Collateral.Add(data.Collateral...))
				if err != nil {
					k.Logger.Error(err.Error())
					return
				}
			} else {
				data.Status = types.Status_STATUS_VERIFIED
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					k.Logger.Error(err.Error())
					return
				}
				publisher := sdk.MustAccAddressFromBech32(data.Publisher)
				err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, publisher, data.Collateral.Add(data.Collateral...))
				if err != nil {
					k.Logger.Error(err.Error())
					return
				}
			}

			// Handle fault validators
			for _, valAddr := range faultValidators {
				k.SetFaultCounter(ctx, valAddr, k.GetFaultCounter(ctx, valAddr)+1)
			}

			// Clean up proofs data
			for _, proof := range proofs {
				addr, err := k.addressCodec.StringToBytes(proof.Sender)
				if err != nil {
					k.Logger.Error(err.Error())
					continue
				}
				k.DeleteProof(sdkCtx, proof.MetadataUri, addr)
			}
		}
	}
}
