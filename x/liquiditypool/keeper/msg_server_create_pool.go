package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := k.addressCodec.StringToBytes(msg.Authority); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// end static validation

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

	pool := types.Pool{
		DenomBase:  msg.DenomBase,
		DenomQuote: msg.DenomQuote,
		FeeRate:    feeRate.String(),
		TickParams: types.TickParams{
			PriceRatio: priceRatio.String(),
			BaseOffset: baseOffset.String(),
		},
		CurrentTick:          0,
		CurrentTickLiquidity: math.LegacyZeroDec().String(),
		CurrentSqrtPrice:     math.LegacyZeroDec().String(),
		Creator:              msg.Creator,
	}

	id, err := k.AppendPool(ctx, pool)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to append pool")
	}

	if err := k.createFeeAccumulator(ctx, id); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&types.EventCreatePool{
		PoolId:     id,
		DenomBase:  msg.DenomBase,
		DenomQuote: msg.DenomQuote,
		FeeRate:    feeRate.String(),
		PriceRatio: priceRatio.String(),
		BaseOffset: baseOffset.String(),
	}); err != nil {
		return nil, err
	}

	return &types.MsgCreatePoolResponse{
		Id: id,
	}, nil
}
