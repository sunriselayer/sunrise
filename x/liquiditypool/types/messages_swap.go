package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwapExactAmountIn{}

func NewMsgSwapExactAmountIn(sender string, tokenIn sdk.Coin, minAmountOut math.Int, routes []SwapRoute) *MsgSwapExactAmountIn {
	return &MsgSwapExactAmountIn{
		Sender:       sender,
		TokenIn:      tokenIn,
		MinAmountOut: minAmountOut,
		Routes:       routes,
	}
}

func (msg *MsgSwapExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if err := msg.TokenIn.Validate(); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid tokenIn (%s)", err)
	}

	if msg.MinAmountOut.LT(math.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid min amount out %s", msg.MinAmountOut.String())
	}

	if len(msg.Routes) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "routes cannot be empty")
	}

	observed := make(map[string]bool)
	for _, route := range msg.Routes {
		if err := ValidateWeights(route.PoolWeights); err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}

		pairName := fmt.Sprintf("%s/%s", route.BaseDenom, route.QuoteDenom)
		if observed[pairName] {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "duplicate pair")
		}
		observed[pairName] = true
	}

	return nil
}

var _ sdk.Msg = &MsgSwapExactAmountOut{}

func NewMsgSwapExactAmountOut(sender string, tokenOut sdk.Coin, maxAmountIn math.Int, route SwapRoute) *MsgSwapExactAmountOut {
	return &MsgSwapExactAmountOut{
		Sender:      sender,
		TokenOut:    tokenOut,
		MaxAmountIn: maxAmountIn,
		Route:       route,
	}
}

func (msg *MsgSwapExactAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if err := msg.TokenOut.Validate(); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid tokenOut (%s)", err)
	}

	if msg.TokenOut.Amount.LTE(math.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid tokenOut amount %s", msg.TokenOut.Amount.String())
	}

	if msg.MaxAmountIn.LTE(math.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid max amount in %s", msg.MaxAmountIn.String())
	}

	if err := ValidateWeights(msg.Route.PoolWeights); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}
