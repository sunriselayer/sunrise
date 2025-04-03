package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/sunriselayer/sunrise/x/lockup/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) NonVotingDelegate(ctx context.Context, msg *types.MsgNonVotingDelegate) (*types.MsgNonVotingDelegateResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Owner); err != nil {
		return nil, errorsmod.Wrap(err, "invalid owner address")
	}

	address := k.LockupAccountAddress(msg.Owner)

	_, err := k.MsgRouterService.Invoke(ctx, &shareclasstypes.MsgNonVotingDelegate{
		Sender:    address.String(),
		Validator: msg.Validator,
		Amount:    msg.Amount,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgNonVotingDelegateResponse{}, nil
}
