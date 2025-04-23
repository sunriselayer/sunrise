package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/lockup/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) NonVotingDelegate(ctx context.Context, msg *types.MsgNonVotingDelegate) (*types.MsgNonVotingDelegateResponse, error) {
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
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "delegate amount denom must be equal to fee denom")
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

	err = k.TrackDelegation(ctx, owner, msg.Id, balance.Amount, lockedAmount, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}

	_, err = k.MsgRouterService.Invoke(ctx, &shareclasstypes.MsgNonVotingDelegate{
		Sender:           lockup.Address,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           msg.Amount,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgNonVotingDelegateResponse{}, nil
}
