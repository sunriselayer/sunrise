package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k msgServer) InitLockupAccount(ctx context.Context, msg *types.MsgInitLockupAccount) (*types.MsgInitLockupAccountResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	owner, err := k.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid owner address")
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Denom != feeDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount denom must be fee denom")
	}

	// Get the current ID for the owner and increment the counter
	id, _, err := k.GetAndIncrementNextLockupAccountID(ctx, owner)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get next lockup account ID")
	}

	// Generate the seed and the lockup account address
	seed := k.makeAddressSeed(owner, id)
	lockupAccAddr := authtypes.NewModuleAddress(seed)

	err = k.bankKeeper.SendCoins(ctx, sender, lockupAccAddr, sdk.NewCoins(msg.Amount))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins")
	}

	lockupAccount := types.LockupAccount{
		Owner:                  msg.Sender,
		Id:                     id,
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
