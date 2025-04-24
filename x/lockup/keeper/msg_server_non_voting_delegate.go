package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k msgServer) NonVotingDelegate(ctx context.Context, msg *types.MsgNonVotingDelegate) (*types.MsgNonVotingDelegateResponse, error) {
	owner, err := k.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid owner address")
	}
	valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}
	lockup, err := k.GetLockupAccount(ctx, owner, msg.Id)
	if err != nil {
		return nil, err
	}
	lockupAddr, err := k.addressCodec.StringToBytes(lockup.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid lockup address")
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Denom != feeDenom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "delegate amount denom must be equal to fee denom")
	}

	balance := k.bankKeeper.GetBalance(ctx, lockupAddr, feeDenom)

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
	_, rewards, err := k.shareclassKeeper.Delegate(ctx, lockupAddr, valAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Add rewards to lockup account
	found, coin := rewards.Find(feeDenom)

	if found {
		err = k.AddRewardsToLockupAccount(ctx, owner, msg.Id, coin.Amount)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgNonVotingDelegateResponse{}, nil
}
