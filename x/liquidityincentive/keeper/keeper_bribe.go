package keeper

import (
	"context"

	"cosmossdk.io/collections"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// SaveVoteWeightsForBribes saves vote weights for bribes distribution
func (k Keeper) SaveVoteWeightsForBribes(ctx context.Context, epochId uint64) error {
	// プールごとの総投票重みを計算
	poolTotalWeights := make(map[uint64]math.LegacyDec)

	// すべての投票を取得して処理
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

	// 各投票者の相対的な重みを保存
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
				// ブライブが存在するプールのみ処理
				bribeKey := collections.Join(epochId, poolWeight.PoolId)
				_, found, err := k.Bribes.Get(ctx, bribeKey)
				if err != nil {
					return false, err
				}
				if !found {
					continue
				}

				// 相対的な重みを計算
				relativeWeight := weight.Quo(poolTotalWeights[poolWeight.PoolId])

				// UnclaimedBribeを保存
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

	// 終了するエポックのIDを取得
	currentEpochId := k.GetCurrentEpochId(ctx)

	// 投票重みを保存
	if err := k.SaveVoteWeightsForBribes(ctx, currentEpochId); err != nil {
		k.Logger(sdk.UnwrapSDKContext(ctx)).Error(
			"failed to save vote weights for bribes",
			"epoch_id", currentEpochId,
			"error", err,
		)
	}

	// 請求期間が終了した古いエポックの未請求ブライブを処理
	if currentEpochId > k.GetParams(ctx).BribeClaimEpochs {
		epochToProcess := currentEpochId - k.GetParams(ctx).BribeClaimEpochs
		if err := k.ProcessUnclaimedBribes(ctx, epochToProcess); err != nil {
			// エラーをログに記録するだけで、エポック終了処理は続行
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

	// エポックが存在するか確認
	_, found, err := k.Epochs.Get(ctx, epochId)
	if err != nil {
		return err
	}
	if !found {
		return errorsmod.Wrapf(types.ErrEpochNotFound, "epoch %d not found", epochId)
	}

	// このエポックのすべてのブライブを処理
	totalUnclaimed := sdk.NewCoins()

	err = k.Bribes.Walk(ctx, collections.NewPrefixedPairRange[uint64, uint64](epochId),
		func(key collections.Pair[uint64, uint64], bribe types.Bribe) (bool, error) {
			// 未請求額を計算
			unclaimedAmount := bribe.Amount.Sub(bribe.ClaimedAmount)
			if !unclaimedAmount.IsZero() {
				totalUnclaimed = totalUnclaimed.Add(unclaimedAmount)
			}

			// ブライブを削除（もう必要ない）
			if err := k.Bribes.Remove(ctx, key); err != nil {
				return false, err
			}

			return false, nil
		})

	if err != nil {
		return err
	}

	// このエポックの未請求ブライブをすべて削除
	err = k.UnclaimedBribes.Walk(ctx, collections.NewPrefixedTripleRange[string, uint64, uint64]("", epochId, 0),
		func(key collections.Triple[string, uint64, uint64], unclaimed types.UnclaimedBribe) (bool, error) {
			return false, k.UnclaimedBribes.Remove(ctx, key)
		})

	if err != nil {
		return err
	}

	// 未請求のブライブがあれば、fee collectorに送信
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

		// イベントを発行
		if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventUnclaimedBribesProcessed{
			EpochId: epochId,
			Amount:  totalUnclaimed,
		}); err != nil {
			return err
		}
	}

	return nil
}
