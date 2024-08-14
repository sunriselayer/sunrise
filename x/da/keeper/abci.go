package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/fee/types"
)

func (k Keeper) EndBlocker(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)
	challengePeriodData, err := k.GetUnverifiedDataBeforeTime(sdkCtx, uint64(sdkCtx.BlockTime().Add(-params.ChallengePeriod).Unix()))
	if err != nil {
		k.Logger().Error(err.Error())
		return
	}

	for _, data := range challengePeriodData {
		if data.Status == "vote_extension" {
			data.Status = "verified"
		}
		if err = k.SetPublishedData(ctx, data); err != nil {
			return
		}

		publisher := sdk.MustAccAddressFromBech32(data.Publisher)
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, publisher, data.Collateral)
		if err != nil {
			k.Logger().Error(err.Error())
			return
		}
	}

	proofPeriodData, err := k.GetUnverifiedDataBeforeTime(sdkCtx, uint64(sdkCtx.BlockTime().Add(-params.ChallengePeriod-params.ProofPeriod).Unix()))
	if err != nil {
		k.Logger().Error(err.Error())
		return
	}

	numActiveValidators := int64(0)
	// votingPowers := make(map[string]int64)
	// powerReduction := k.StakingKeeper.PowerReduction(ctx)
	iterator, err := k.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		k.Logger().Error(err.Error())
		return
	}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validator, err := k.StakingKeeper.Validator(ctx, iterator.Value())
		if err != nil {
			k.Logger().Error(err.Error())
			return
		}

		if validator.IsBonded() {
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

	for _, data := range proofPeriodData {
		if data.Status == "challenge_for_fraud" {
			// bondedTokens, err := k.StakingKeeper.TotalBondedTokens(ctx)
			// if err != nil {
			// 	k.Logger().Error(err.Error())
			// 	return
			// }

			// totalBondedPower := sdk.TokensToConsensusPower(bondedTokens, k.StakingKeeper.PowerReduction(ctx))
			// thresholdPower := params.VoteThreshold.MulInt64(totalBondedPower).RoundInt().Int64()
			proofs := k.GetProofs(sdkCtx, data.MetadataUri)
			shardProofCount := make(map[int64]int64)
			for _, proof := range proofs {
				for _, indice := range proof.Indices {
					shardProofCount[indice]++
				}
			}

			safeShardCount := int64(0)
			for _, proofCount := range shardProofCount {
				// len(zkp_including_this_shard) / replication_factor >= 2/3
				if math.LegacyNewDec(proofCount).GTE(params.ReplicationFactor.MulInt64(2).QuoInt64(3)) {
					safeShardCount++
				}
			}

			// valid_shards / len(shards) >= 1/2
			if safeShardCount*2 < int64(len(data.ShardDoubleHashes)) {
				// TODO: might require rejected records as well
				data.Status = "rejected"
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					k.Logger().Error(err.Error())
					return
				}
				// k.DeletePublishedData(sdkCtx, data)

				challenger := sdk.MustAccAddressFromBech32(data.Challenger)
				err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, challenger, data.Collateral.Add(data.Collateral...))
				if err != nil {
					k.Logger().Error(err.Error())
					return
				}
			} else {
				data.Status = "verified"
				err = k.SetPublishedData(ctx, data)
				if err != nil {
					k.Logger().Error(err.Error())
					return
				}
				publisher := sdk.MustAccAddressFromBech32(data.Publisher)
				err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, publisher, data.Collateral.Add(data.Collateral...))
				if err != nil {
					k.Logger().Error(err.Error())
					return
				}
			}

			// TODO: count validators not voted in proof submission phase
			// Clean up proofs data
			for _, proof := range proofs {
				k.DeleteProof(sdkCtx, proof.MetadataUri, proof.Sender)
			}
		}
	}
}
