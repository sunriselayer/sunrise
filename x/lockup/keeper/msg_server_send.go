package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k msgServer) Send(ctx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	owner, err := k.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid owner address")
	}
	recipient, err := k.addressCodec.StringToBytes(msg.Recipient)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid recipient address")
	}

	lockupAccount, err := k.GetLockupAccount(ctx, owner)
	if err != nil {
		return nil, err
	}

	address := k.LockupAccountAddress(msg.Owner)

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	found, coin := msg.Amount.Find(feeDenom)

	if found {
		totalLockupAmount := lockupAccount.LockupAmountOriginal.Add(lockupAccount.LockupAmountAdditional)
		balance := k.bankKeeper.GetBalance(ctx, address, feeDenom)
		sendAmount := coin.Amount

		sdkCtx := sdk.UnwrapSDKContext(ctx)
		now := sdkCtx.BlockTime()

		unlockedAmount, err := types.CalculateUnlockedAmount(totalLockupAmount, lockupAccount.StartTime, lockupAccount.EndTime, now)
		if err != nil {
			return nil, err
		}

		canSend := types.SendCondition(totalLockupAmount, unlockedAmount, balance.Amount, sendAmount)
		if !canSend {
			return nil, errorsmod.Wrap(types.ErrInsufficientUnlockedFunds, "insufficient unlocked funds")
		}
	}

	err = k.bankKeeper.SendCoins(ctx, address, recipient, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendResponse{}, nil
}
