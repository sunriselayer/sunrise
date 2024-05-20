package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Authority != k.authority {
		return nil, types.ErrInvalidSigner
	}

	var pool = types.Pool{
		Id:                   0,
		DenomBase:            msg.DenomBase,
		DenomQuote:           msg.DenomQuote,
		FeeRate:              msg.FeeRate,
		TickParams:           types.TickParams{}, // TODO:
		CurrentTick:          0,
		CurrentTickLiquidity: math.LegacyZeroDec(),
		CurrentSqrtPrice:     math.LegacyZeroDec(),
	}
	id := k.AppendPool(ctx, pool)

	return &types.MsgCreatePoolResponse{
		Id: id,
	}, nil
}
