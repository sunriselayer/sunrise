package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k msgServer) NonVotingUndelegate(ctx context.Context, msg *types.MsgNonVotingUndelegate) (*types.MsgNonVotingUndelegateResponse, error) {
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

	output, rewards, completionTime, err := k.shareclassKeeper.Undelegate(ctx, lockupAddr, lockupAddr, valAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	isNewEntry := true
	if lockup.UnbondEntries == nil {
		lockup.UnbondEntries = &types.UnbondingEntries{}
	}
	entries := lockup.UnbondEntries.Entries
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	for i, entry := range entries {
		if entry.CreationHeight == sdkCtx.BlockHeight() && entry.EndTime == completionTime.Unix() {
			entry.Amount = entry.Amount.Add(output.Amount)

			// update the entry
			entries[i] = entry
			isNewEntry = false
			break
		}
	}

	if isNewEntry {
		entries = append(entries, &types.UnbondingEntry{
			CreationHeight:   sdkCtx.BlockHeight(),
			EndTime:          completionTime.Unix(),
			Amount:           output.Amount,
			ValidatorAddress: msg.ValidatorAddress,
		})
	}

	lockup.UnbondEntries.Entries = entries
	err = k.SetLockupAccount(ctx, lockup)
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

	return &types.MsgNonVotingUndelegateResponse{}, nil
}
