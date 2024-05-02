package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgJoinPool{}

func NewMsgJoinPool(sender string, poolId uint64, baseToken sdk.Coin, quoteToken sdk.Coin, minShareAmount math.Int) *MsgJoinPool {
	return &MsgJoinPool{
		Sender:         sender,
		PoolId:         poolId,
		BaseToken:      baseToken,
		QuoteToken:     quoteToken,
		MinShareAmount: minShareAmount,
	}
}

func (msg *MsgJoinPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if err := msg.BaseToken.Validate(); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid base token (%s)", err)
	}

	if err := msg.QuoteToken.Validate(); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid quote token (%s)", err)
	}

	if msg.MinShareAmount.LT(math.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid min share amount %s", msg.MinShareAmount.String())
	}

	return nil
}

var _ sdk.Msg = &MsgExitPool{}

func NewMsgExitPool(sender string, poolId uint64, shareAmount math.Int, minAmountBase math.Int, minAmountQuote math.Int) *MsgExitPool {
	return &MsgExitPool{
		Sender:         sender,
		PoolId:         poolId,
		ShareAmount:    shareAmount,
		MinAmountBase:  minAmountBase,
		MinAmountQuote: minAmountQuote,
	}
}

func (msg *MsgExitPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.ShareAmount.LTE(math.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid share amount %s", msg.ShareAmount.String())
	}

	if msg.MinAmountBase.LT(math.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid min amount base %s", msg.MinAmountBase.String())
	}

	if msg.MinAmountQuote.LT(math.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid min amount quote %s", msg.MinAmountQuote.String())
	}

	return nil
}
