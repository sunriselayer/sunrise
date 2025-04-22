package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/sunriselayer/sunrise/x/lockup/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k msgServer) NonVotingUndelegate(ctx context.Context, msg *types.MsgNonVotingUndelegate) (*types.MsgNonVotingUndelegateResponse, error) {
	owner, err := k.addressCodec.StringToBytes(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid owner address")
	}

	lockupAcc := k.LockupAccountAddress(owner, msg.Id)

	_, err = k.MsgRouterService.Invoke(ctx, &shareclasstypes.MsgNonVotingUndelegate{
		Sender:           lockupAcc.String(),
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           msg.Amount,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgNonVotingUndelegateResponse{}, nil
}
