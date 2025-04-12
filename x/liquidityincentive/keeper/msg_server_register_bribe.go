package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k msgServer) RegisterBribe(ctx context.Context, msg *types.MsgRegisterBribe) (*types.MsgRegisterBribeResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	// エポックが存在するか確認
	epoch, found, err := k.Epochs.Get(ctx, msg.EpochId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrEpochNotFound, "epoch %d not found", msg.EpochId)
	}

	// プールが存在するか確認
	_, found, err = k.liquidityPoolKeeper.GetPool(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrPoolNotFound, "pool %d not found", msg.PoolId)
	}

	// ブライブの金額が有効か確認
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return nil, errorsmod.Wrap(types.ErrInvalidBribeAmount, "bribe amount must be valid and non-zero")
	}

	// 送信者からブライブの金額を引き出す
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(sdkCtx, sender, types.ModuleName, msg.Amount); err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to module")
	}

	// ブライブを保存または更新
	key := collections.Join(msg.EpochId, msg.PoolId)
	existingBribe, found, err := k.Bribes.Get(ctx, key)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get existing bribe")
	}

	if found {
		// 既存のブライブに追加
		existingBribe.Amount = existingBribe.Amount.Add(msg.Amount)
		if err := k.Bribes.Set(ctx, key, existingBribe); err != nil {
			return nil, errorsmod.Wrap(err, "failed to update bribe")
		}
	} else {
		// 新しいブライブを作成
		bribe := types.Bribe{
			EpochId:       msg.EpochId,
			PoolId:        msg.PoolId,
			Amount:        msg.Amount,
			ClaimedAmount: sdk.NewCoin(msg.Amount.Denom, sdk.ZeroInt()),
		}

		if err := k.Bribes.Set(ctx, key, bribe); err != nil {
			return nil, errorsmod.Wrap(err, "failed to set bribe")
		}
	}

	// イベントを発行
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventRegisterBribe{
		Sender:  msg.Sender,
		EpochId: msg.EpochId,
		PoolId:  msg.PoolId,
		Amount:  msg.Amount,
	}); err != nil {
		return nil, err
	}

	return &types.MsgRegisterBribeResponse{}, nil
}
