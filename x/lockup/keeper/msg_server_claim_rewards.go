package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
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

	res, err := k.ShareclassMsgServer.ClaimRewards(ctx, &shareclasstypes.MsgClaimRewards{
		Sender:           lockup.Address,
		ValidatorAddress: msg.ValidatorAddress,
	})
	if err != nil {
		return nil, err
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	found, coin := res.Amount.Find(feeDenom)

	if found {
		err = k.AddRewardsToLockupAccount(ctx, owner, msg.Id, coin.Amount)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgClaimRewardsResponse{}, nil
}
