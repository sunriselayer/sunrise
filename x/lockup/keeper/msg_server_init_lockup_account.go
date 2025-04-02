package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k msgServer) InitLockupAccount(ctx context.Context, msg *types.MsgInitLockupAccount) (*types.MsgInitLockupAccountResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, err
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Denom != feeDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount denom must be fee denom")
	}

	_, err = k.GetLockupAccount(ctx, sender)
	if err == nil {
		return nil, errorsmod.Wrap(types.ErrLockupAccountAlreadyExists, "lockup account already exists")
	}

	lockupAccount := types.LockupAccount{
		Owner:                  msg.Sender,
		StartTime:              msg.StartTime,
		EndTime:                msg.EndTime,
		LockupAmountOriginal:   msg.Amount.Amount,
		LockupAmountAdditional: math.ZeroInt(),
	}

	err = k.SetLockupAccount(ctx, lockupAccount)
	if err != nil {
		return nil, err
	}

	return &types.MsgInitLockupAccountResponse{}, nil
}
