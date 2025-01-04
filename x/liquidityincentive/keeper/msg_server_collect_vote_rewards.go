package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CollectVoteRewards(ctx context.Context, msg *types.MsgCollectVoteRewards) (*types.MsgCollectVoteRewardsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// Currently BeginBlocker distributes a portion of vRISE to the gauge voters
	// In the future, a portion of swap fee rewards will be distributed to the gauge voters
	implemented := false
	if !implemented {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotSupported, "This message is reserved for the future implementation and not supported yet")
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCollectVoteRewards{
		Address: msg.Sender,
	}); err != nil {
		return nil, err
	}

	return &types.MsgCollectVoteRewardsResponse{}, nil
}
