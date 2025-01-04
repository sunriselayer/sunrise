package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func (q queryServer) CalculationSwapExactAmountIn(goCtx context.Context, req *types.QueryCalculationSwapExactAmountInRequest) (*types.QueryCalculationSwapExactAmountInResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	amountIn, ok := sdkmath.NewIntFromString(req.AmountIn)
	if !ok {
		return nil, types.ErrInvalidAmount
	}
	result, interfaceProviderFee, err := q.k.CalculateResultExactAmountIn(ctx, req.HasInterfaceFee, *req.Route, amountIn)
	if err != nil {
		return nil, err
	}

	return &types.QueryCalculationSwapExactAmountInResponse{
		Result:               result,
		InterfaceProviderFee: interfaceProviderFee,
		AmountOut:            result.TokenOut.Amount.Sub(interfaceProviderFee),
	}, nil
}

func (q queryServer) CalculationSwapExactAmountOut(goCtx context.Context, req *types.QueryCalculationSwapExactAmountOutRequest) (*types.QueryCalculationSwapExactAmountOutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	amountOut, ok := sdkmath.NewIntFromString(req.AmountOut)
	if !ok {
		return nil, types.ErrInvalidAmount
	}
	result, interfaceProviderFee, err := q.k.CalculateResultExactAmountOut(ctx, req.HasInterfaceFee, *req.Route, amountOut)
	if err != nil {
		return nil, err
	}

	return &types.QueryCalculationSwapExactAmountOutResponse{
		Result:               result,
		InterfaceProviderFee: interfaceProviderFee,
		AmountIn:             result.TokenIn.Amount,
	}, nil
}
