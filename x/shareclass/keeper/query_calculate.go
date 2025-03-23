package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (q queryServer) CalculateAmount(ctx context.Context, req *types.QueryCalculateAmountRequest) (*types.QueryCalculateAmountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	params, err := q.k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	amount, err := q.k.CalculateAmountByShare(ctx, req.ValidatorAddress, req.Share)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	coin := sdk.NewCoin(params.BondDenom, amount)

	return &types.QueryCalculateAmountResponse{Amount: coin}, nil
}
