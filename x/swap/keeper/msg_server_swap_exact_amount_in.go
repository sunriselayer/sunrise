package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/swap/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) SwapExactAmountIn(ctx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
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

	if !msg.AmountIn.IsPositive() {
		return nil, errorsmod.Wrapf(types.ErrInvalidAmount, "amount in must be positive: %s", msg.AmountIn)
	}

	if !msg.MinAmountOut.IsPositive() {
		return nil, errorsmod.Wrapf(types.ErrInvalidAmount, "min amount out must be positive: %s", msg.MinAmountOut)
	}
	// end static validation

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	result, interfaceProviderFee, err := k.Keeper.SwapExactAmountIn(sdkCtx, sender, msg.InterfaceProvider, msg.Route, msg.AmountIn, msg.MinAmountOut)
	if err != nil {
		return nil, err
	}

	return &types.MsgSwapExactAmountInResponse{
		Result:               result,
		InterfaceProviderFee: interfaceProviderFee,
		AmountOut:            result.TokenOut.Amount.Sub(interfaceProviderFee),
	}, nil
}
