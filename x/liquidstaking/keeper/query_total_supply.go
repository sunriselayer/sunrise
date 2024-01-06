package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunrise-zone/sunrise-app/x/liquidstaking/types"
)

func (k Keeper) TotalSupply(goCtx context.Context, req *types.QueryTotalSupplyRequest) (*types.QueryTotalSupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	totalValue, err := k.GetTotalDerivativeValue(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryTotalSupplyResponse{
		Height: ctx.BlockHeight(),
		Result: []sdk.Coin{totalValue},
	}, nil
}
