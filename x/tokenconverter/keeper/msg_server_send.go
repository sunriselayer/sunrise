package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k msgServer) Send(ctx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	fromAddress, err := k.addressCodec.StringToBytes(msg.FromAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid from address")
	}

	toAddress, err := k.addressCodec.StringToBytes(msg.ToAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid to address")
	}

	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be positive")
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	isAllowed := false
	for _, addr := range params.AllowedAddresses {
		if addr == msg.FromAddress || addr == msg.ToAddress {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "from or to address is not in allowed addresses")
	}

	if err := k.bankKeeper.SendCoins(ctx, sdk.AccAddress(fromAddress), sdk.AccAddress(toAddress), msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgSendResponse{}, nil
}
