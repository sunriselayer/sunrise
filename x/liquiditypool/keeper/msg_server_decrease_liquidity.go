package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DecreaseLiquidity(ctx context.Context, msg *types.MsgDecreaseLiquidity) (*types.MsgDecreaseLiquidityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	// end static validation

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	liquidity, err := math.LegacyNewDecFromStr(msg.Liquidity)
	if err != nil {
		return nil, err
	}
	amountBase, amountQuote, err := k.Keeper.DecreaseLiquidity(sdkCtx, sender, msg.Id, liquidity)
	if err != nil {
		return nil, err
	}

	return &types.MsgDecreaseLiquidityResponse{
		AmountBase:  amountBase,
		AmountQuote: amountQuote,
	}, nil
}
