package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: enable later after completion of module - for now, disabled for testing
	// if msg.Authority != k.authority {
	// 	return nil, types.ErrInvalidSigner
	// }

	feeRate, err := math.LegacyNewDecFromStr(msg.FeeRate)
	if err != nil {
		return nil, err
	}

	priceRatio, err := math.LegacyNewDecFromStr(msg.PriceRatio)
	if err != nil {
		return nil, err
	}

	baseOffset, err := math.LegacyNewDecFromStr(msg.BaseOffset)
	if err != nil {
		return nil, err
	}

	var pool = types.Pool{
		Id:         0,
		DenomBase:  msg.DenomBase,
		DenomQuote: msg.DenomQuote,
		FeeRate:    feeRate,
		TickParams: types.TickParams{
			PriceRatio: priceRatio,
			BaseOffset: baseOffset,
		},
		CurrentTick:          0,
		CurrentTickLiquidity: math.LegacyZeroDec(),
		CurrentSqrtPrice:     math.LegacyZeroDec(),
	}
	id := k.AppendPool(ctx, pool)
	if err := k.createFeeAccumulator(ctx, id); err != nil {
		return nil, err
	}

	return &types.MsgCreatePoolResponse{
		Id: id,
	}, nil
}
