package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CollectVoteRewards(goCtx context.Context, msg *types.MsgCollectVoteRewards) (*types.MsgCollectVoteRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

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
