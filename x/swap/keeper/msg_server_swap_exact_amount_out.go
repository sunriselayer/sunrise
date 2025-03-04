package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/swap/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) SwapExactAmountOut(ctx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	if msg.InterfaceProvider != "" {
		if _, err := sdk.AccAddressFromBech32(msg.InterfaceProvider); err != nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid interface provider address (%s)", err)
		}
	}

	if err := msg.Route.Validate(); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidRoute, "invalid route: %s", err)
	}

	if !msg.MaxAmountIn.IsPositive() {
		return nil, errorsmod.Wrapf(types.ErrInvalidAmount, "max amount in must be positive: %s", msg.MaxAmountIn)
	}

	if !msg.AmountOut.IsPositive() {
		return nil, errorsmod.Wrapf(types.ErrInvalidAmount, "amount out must be positive: %s", msg.AmountOut)
	}
	// end static validation

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	result, interfaceProviderFee, err := k.Keeper.SwapExactAmountOut(sdkCtx, sender, msg.InterfaceProvider, msg.Route, msg.MaxAmountIn, msg.AmountOut)
	if err != nil {
		return nil, err
	}

	return &types.MsgSwapExactAmountOutResponse{
		Result:               result,
		InterfaceProviderFee: interfaceProviderFee,
		AmountOut:            result.TokenOut.Amount.Sub(interfaceProviderFee),
	}, nil
}
