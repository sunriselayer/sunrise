package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePool(ctx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Authority); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	err := sdk.ValidateDenom(msg.DenomBase)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid denom base")
	}

	err = sdk.ValidateDenom(msg.DenomQuote)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid denom quote")
	}

	feeRate, err := math.LegacyNewDecFromStr(msg.FeeRate)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid fee rate")
	}

	if !feeRate.IsPositive() {
		return nil, errorsmod.Wrap(err, "fee rate must be positive")
	}

	if feeRate.GTE(math.LegacyOneDec()) {
		return nil, errorsmod.Wrap(err, "fee rate must be less than 1")
	}

	if msg.PriceRatio != "1.0001" {
		return nil, errorsmod.Wrap(err, "price ratio must be 1.0001")
	}

	priceRatio, err := math.LegacyNewDecFromStr(msg.PriceRatio)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid price ratio")
	}

	baseOffset, err := math.LegacyNewDecFromStr(msg.BaseOffset)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid base offset")
	}

	if baseOffset.IsPositive() {
		return nil, errorsmod.Wrap(err, "base offset must be in range (-1, 0]")
	}

	if baseOffset.LTE(math.LegacyNewDec(-1)) {
		return nil, errorsmod.Wrap(err, "base offset must be in range (-1, 0]")
	}

	// end static validation

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get params")
	}

	if params.CreatePoolGas > 0 {
		sdk.UnwrapSDKContext(ctx).GasMeter().ConsumeGas(params.CreatePoolGas, "create pool")
	}

	var pool = types.Pool{
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
	}
	id, err := k.AppendPool(ctx, pool)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to append pool")
	}
	if err := k.createFeeAccumulator(ctx, id); err != nil {
		return nil, errorsmod.Wrap(err, "failed to create fee accumulator")
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreatePool{
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
