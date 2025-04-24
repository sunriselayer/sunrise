package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/lockup/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) NonVotingUndelegate(ctx context.Context, msg *types.MsgNonVotingUndelegate) (*types.MsgNonVotingUndelegateResponse, error) {
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
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "undelegate amount denom must be equal to fee denom")
	}

	res, err := k.MsgRouterService.Invoke(ctx, &shareclasstypes.MsgNonVotingUndelegate{
		Sender:           lockup.Address,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           msg.Amount,
	})
	if err != nil {
		return nil, err
	}

	undelegateResponse, ok := res.(*shareclasstypes.MsgNonVotingUndelegateResponse)
	if !ok {
		return nil, sdkerrors.ErrInvalidRequest
	}
	if undelegateResponse == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "undelegate response is nil")
	}

	isNewEntry := true
	entries := lockup.UnbondEntries.Entries
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	for i, entry := range entries {
		if entry.CreationHeight == sdkCtx.BlockHeight() && entry.EndTime == undelegateResponse.CompletionTime.Unix() {
			entry.Amount = entry.Amount.Add(undelegateResponse.Amount.Amount)

			// update the entry
			entries[i] = entry
			isNewEntry = false
			break
		}
	}

	if isNewEntry {
		entries = append(entries, &types.UnbondingEntry{
			CreationHeight:   sdkCtx.BlockHeight(),
			EndTime:          undelegateResponse.CompletionTime.Unix(),
			Amount:           undelegateResponse.Amount.Amount,
			ValidatorAddress: msg.ValidatorAddress,
		})
	}

	lockup.UnbondEntries.Entries = entries
	err = k.SetLockupAccount(ctx, lockup)
	if err != nil {
		return nil, err
	}

	// Add rewards to lockup account
	found, coin := undelegateResponse.Rewards.Find(feeDenom)
	if found {
		err = k.AddRewardsToLockupAccount(ctx, owner, msg.Id, coin.Amount)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgNonVotingUndelegateResponse{}, nil
}
