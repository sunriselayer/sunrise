package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) IncentivesByPositionIds(goCtx context.Context, req *types.QueryIncentivesByPositionIdsRequest) (*types.QueryIncentivesByPositionIdsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	// TODO: implement

	return &types.QueryIncentivesByPositionIdsResponse{}, nil
}
