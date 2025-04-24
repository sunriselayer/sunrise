package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	banktypes "cosmossdk.io/x/bank/types"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k msgServer) Send(ctx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	owner, err := k.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid owner address")
	}

	lockup, err := k.GetLockupAccount(ctx, owner, msg.Id)
	if err != nil {
		return nil, err
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Denom != feeDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount denom must be equal to fee denom")
	}

	lockupAcc, err := k.addressCodec.StringToBytes(lockup.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid lockup account address")
	}

	balance := k.bankKeeper.GetBalance(ctx, lockupAcc, feeDenom)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime().Unix()

	_, lockedAmount, err := lockup.GetLockCoinInfo(currentTime)
	if err != nil {
		return nil, err
	}

	// refresh ubd entries to make sure delegation locking amount is up to date
	err = k.CheckUnbondingEntriesMature(ctx, owner, msg.Id)
	if err != nil {
		return nil, err
	}

	notBondedLockedAmount := lockup.GetNotBondedLockedAmount(lockedAmount)

	spendable, err := balance.Amount.SafeSub(notBondedLockedAmount)
	if err != nil {
		return nil, errorsmod.Wrapf(err,
			"locked amount exceeds account balance funds: %d > %d", lockedAmount, balance.Amount)
	}

	if spendable.LT(msg.Amount.Amount) {
		return nil, errorsmod.Wrapf(err,
			"spendable balance %d is smaller than %d",
			spendable, msg.Amount.Amount,
		)
	}

	_, err = k.MsgRouterService.Invoke(ctx, &banktypes.MsgSend{
		FromAddress: lockup.Address,
		ToAddress:   msg.Recipient,
		Amount:      sdk.NewCoins(msg.Amount),
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgSendResponse{}, nil
}
