package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"errors"

	"cosmossdk.io/collections"
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
	lockup, err := k.GetLockupAccount(ctx, owner, msg.LockupAccountId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrapf(types.ErrLockupAccountNotFound, "owner %s, id %d", msg.Owner, msg.LockupAccountId)
		}
		return nil, err
	}
	err = msg.Amount.Validate()
	if err != nil {
		return nil, err
	}

	lockupAddr, err := k.addressCodec.StringToBytes(lockup.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid lockup address")
	}

	transferableDenom, err := k.tokenConverterKeeper.GetTransferableDenom(ctx)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Denom != transferableDenom {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid denom: expected %s, got %s", transferableDenom, msg.Amount.Denom)
	}

	balance := k.bankKeeper.GetBalance(ctx, lockupAddr, transferableDenom)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime().Unix()

	_, locked, err := lockup.GetLockCoinInfo(currentTime)
	if err != nil {
		return nil, err
	}

	// refresh ubd entries to make sure delegation locking amount is up to date
	err = k.CheckUnbondingEntriesMature(ctx, owner, msg.LockupAccountId)
	if err != nil {
		return nil, err
	}

	// Track the delegation
	err = k.TrackDelegation(ctx, owner, msg.LockupAccountId, balance.Amount, locked, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}

	// Delegate the tokens to the validator
	_, rewards, err := k.shareclassKeeper.Delegate(ctx, lockupAddr, valAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Add rewards to lockup account
	found, coin := rewards.Find(transferableDenom)

	if found {
		err = k.AddAdditionalLockup(ctx, owner, msg.LockupAccountId, coin.Amount)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgNonVotingDelegateResponse{}, nil
}
