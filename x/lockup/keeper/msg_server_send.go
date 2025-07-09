package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"errors"

	"cosmossdk.io/collections"
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
	lockupAcc, err := k.GetLockupAccount(ctx, owner, msg.LockupAccountId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, types.ErrLockupAccountNotFound
		}
		return nil, err
	}
	err = msg.Amount.Validate()
	if err != nil {
		return nil, err
	}

	lockupAddr, err := k.addressCodec.StringToBytes(lockupAcc.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid lockup account address")
	}

	transferableDenom, err := k.tokenConverterKeeper.GetTransferableDenom(ctx)
	if err != nil {
		return nil, err
	}

	found, transferableCoin := msg.Amount.Find(transferableDenom)

	// if amount is fee denom, check if the balance is enough
	if found {
		balance := k.bankKeeper.GetBalance(ctx, lockupAddr, transferableDenom)

		sdkCtx := sdk.UnwrapSDKContext(ctx)
		currentTime := sdkCtx.BlockTime().Unix()

		_, lockedAmount, err := lockupAcc.GetLockCoinInfo(currentTime)
		if err != nil {
			return nil, err
		}

		// refresh ubd entries to make sure delegation locking amount is up to date
		err = k.CheckUnbondingEntriesMature(ctx, owner, msg.LockupAccountId)
		if err != nil {
			return nil, err
		}

		notBondedLockedAmount := lockupAcc.GetNotBondedLockedAmount(lockedAmount)

		spendable, err := balance.Amount.SafeSub(notBondedLockedAmount)
		if err != nil {
			return nil, errorsmod.Wrapf(err,
				"locked amount exceeds account balance funds: %d > %d", lockedAmount, balance.Amount)
		}

		if spendable.LT(transferableCoin.Amount) {
			return nil, errorsmod.Wrapf(types.ErrNotEnoughSpendableBalance,
				"spendable balance %s is smaller than %s",
				sdk.NewCoin(transferableDenom, spendable), transferableCoin,
			)
		}
	}

	err = k.bankKeeper.SendCoins(ctx, lockupAddr, recipient, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendResponse{}, nil
}
