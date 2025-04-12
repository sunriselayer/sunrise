package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k msgServer) ClaimBribe(ctx context.Context, msg *types.MsgClaimBribe) (*types.MsgClaimBribeResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	senderAddr := sdk.AccAddress(sender)
	totalClaimed := sdk.NewCoins()

	// エポックが存在するか確認
	_, found, err := k.Epochs.Get(ctx, msg.EpochId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrEpochNotFound, "epoch %d not found", msg.EpochId)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 各プールのブライブを処理
	for _, poolId := range msg.PoolIds {
		// ブライブが存在するか確認
		bribeKey := collections.Join(msg.EpochId, poolId)
		bribe, found, err := k.Bribes.Get(ctx, bribeKey)
		if err != nil {
			return nil, err
		}
		if !found {
			continue // このプールにブライブがない
		}

		// 未請求のブライブを取得
		unclaimedKey := collections.Join3(msg.Sender, msg.EpochId, poolId)
		unclaimed, found, err := k.UnclaimedBribes.Get(ctx, unclaimedKey)
		if err != nil {
			return nil, err
		}
		if !found {
			continue // 請求権がない
		}

		// 重みを取得
		weight, err := math.LegacyNewDecFromStr(unclaimed.Weight)
		if err != nil {
			return nil, errorsmod.Wrap(err, "invalid weight format")
		}

		// 請求額を計算
		claimAmount := sdk.NewCoin(
			bribe.Amount.Denom,
			math.LegacyNewDecFromInt(bribe.Amount.Amount).Mul(weight).TruncateInt(),
		)

		if claimAmount.IsZero() {
			continue
		}

		// ブライブを送信
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			sdkCtx,
			types.ModuleName,
			senderAddr,
			sdk.NewCoins(claimAmount),
		); err != nil {
			return nil, errorsmod.Wrap(err, "failed to send coins from module")
		}

		// 請求済み金額を更新
		bribe.ClaimedAmount = bribe.ClaimedAmount.Add(claimAmount)
		if err := k.Bribes.Set(ctx, bribeKey, bribe); err != nil {
			return nil, errorsmod.Wrap(err, "failed to update bribe claimed amount")
		}

		// UnclaimedBribeを削除（重複請求を防止）
		if err := k.UnclaimedBribes.Remove(ctx, unclaimedKey); err != nil {
			return nil, errorsmod.Wrap(err, "failed to remove unclaimed bribe")
		}

		totalClaimed = totalClaimed.Add(claimAmount)
	}

	if totalClaimed.IsZero() {
		return nil, errorsmod.Wrap(types.ErrNoBribesToClaim, "no bribes to claim")
	}

	// イベントを発行
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventClaimBribe{
		Sender:  msg.Sender,
		EpochId: msg.EpochId,
		PoolIds: msg.PoolIds,
		Amount:  totalClaimed,
	}); err != nil {
		return nil, err
	}

	return &types.MsgClaimBribeResponse{
		ClaimedAmount: totalClaimed,
	}, nil
}
