package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/sunriselayer/sunrise/x/lockup/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) ClaimRewards(ctx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Owner); err != nil {
		return nil, errorsmod.Wrap(err, "invalid owner address")
	}

	address := k.LockupAccountAddress(msg.Owner)

	response, err := k.shareclassMsgServer.ClaimRewards(ctx, &shareclasstypes.MsgClaimRewards{
		Sender:    address.String(),
		Validator: msg.Validator,
	})
	if err != nil {
		return nil, err
	}

	feeDenom, err := k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, err
	}

	found, coin := response.Amount.Find(feeDenom)

	if found {
		lockupAccount, err := k.GetLockupAccount(ctx, address)
		if err != nil {
			return nil, err
		}

		lockupAccount.LockupAmountAdditional = lockupAccount.LockupAmountAdditional.Add(coin.Amount)

		err = k.SetLockupAccount(ctx, lockupAccount)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgClaimRewardsResponse{}, nil
}
