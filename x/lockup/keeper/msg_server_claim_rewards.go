package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/lockup/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) ClaimRewards(ctx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	owner, err := k.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid owner address")
	}

	lockup, err := k.GetLockupAccount(ctx, owner, msg.Id)
	if err != nil {
		return nil, err
	}

	res, err := k.MsgRouterService.Invoke(ctx, &shareclasstypes.MsgClaimRewards{
		Sender:           lockup.Address,
		ValidatorAddress: msg.ValidatorAddress,
	})
	if err != nil {
		return nil, err
	}

	claimRewardsResponse, ok := res.(*shareclasstypes.MsgClaimRewardsResponse)
	if !ok {
		return nil, sdkerrors.ErrInvalidRequest
	}
	if claimRewardsResponse == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid response")
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	found, coin := claimRewardsResponse.Amount.Find(feeDenom)

	if found {
		lockupAccount, err := k.GetLockupAccount(ctx, owner, msg.Id)
		if err != nil {
			return nil, err
		}

		lockupAccount.AdditionalLocking = lockupAccount.AdditionalLocking.Add(coin.Amount)

		err = k.SetLockupAccount(ctx, lockupAccount)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgClaimRewardsResponse{}, nil
}
