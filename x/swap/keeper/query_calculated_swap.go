package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func (k Keeper) CalculatedSwapExactAmountIn(goCtx context.Context, req *types.QueryCalculatedSwapExactAmountInRequest) (*types.QueryCalculatedSwapExactAmountInResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	result, interfaceProviderFee, err := k.CalculateResultExactAmountIn(ctx, req.HasInterfaceFee, req.Route, req.AmountIn)
	if err != nil {
		return nil, err
	}

	return &types.QueryCalculatedSwapExactAmountInResponse{
		Result:               result,
		InterfaceProviderFee: interfaceProviderFee,
		AmountOut:            result.TokenOut.Amount.Sub(interfaceProviderFee),
	}, nil
}

func (k Keeper) CalculatedSwapExactAmountOut(goCtx context.Context, req *types.QueryCalculatedSwapExactAmountOutRequest) (*types.QueryCalculatedSwapExactAmountOutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	result, interfaceProviderFee, err := k.CalculateResultExactAmountOut(ctx, req.HasInterfaceFee, req.Route, req.AmountOut)
	if err != nil {
		return nil, err
	}

	return &types.QueryCalculatedSwapExactAmountOutResponse{
		Result:               result,
		InterfaceProviderFee: interfaceProviderFee,
		AmountIn:             result.TokenIn.Amount,
	}, nil
}
