package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k msgServer) WithdrawUnbonded(ctx context.Context, msg *types.MsgWithdrawUnbonded) (*types.MsgWithdrawUnbondedResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	sender = sender

	return &types.MsgWithdrawUnbondedResponse{}, nil
}
