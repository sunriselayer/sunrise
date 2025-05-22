package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (q queryServer) CalculateBondingAmount(ctx context.Context, req *types.QueryCalculateBondingAmountRequest) (*types.QueryCalculateBondingAmountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	share, ok := sdkmath.NewIntFromString(req.Share)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid share")
	}
	amount, err := q.k.CalculateAmountByShare(ctx, req.ValidatorAddress, share)
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

	amount, ok := sdkmath.NewIntFromString(req.Amount)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}
	share, err := q.k.CalculateShareByAmount(ctx, req.ValidatorAddress, amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCalculateShareResponse{Share: share}, nil
}
