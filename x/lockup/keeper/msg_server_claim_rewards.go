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

	_, err := k.MsgRouterService.Invoke(ctx, &shareclasstypes.MsgClaimRewards{
		Sender:    address.String(),
		Validator: msg.Validator,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimRewardsResponse{}, nil
}
