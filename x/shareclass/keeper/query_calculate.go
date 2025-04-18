package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (q queryServer) CalculateBondingAmount(ctx context.Context, req *types.QueryCalculateBondingAmountRequest) (*types.QueryCalculateBondingAmountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	amount, err := q.k.CalculateAmountByShare(ctx, req.ValidatorAddress, req.Share)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	feeDenom, err := q.k.feeKeeper.FeeDenom(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	coin := sdk.NewCoin(feeDenom, amount)

	return &types.QueryCalculateBondingAmountResponse{Amount: coin}, nil
}

func (q queryServer) CalculateShare(ctx context.Context, req *types.QueryCalculateShareRequest) (*types.QueryCalculateShareResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	share, err := q.k.CalculateShareByAmount(ctx, req.ValidatorAddress, req.Amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCalculateShareResponse{Share: share}, nil
}
