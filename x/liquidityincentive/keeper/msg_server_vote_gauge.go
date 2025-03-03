package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"

	errorsmod "cosmossdk.io/errors"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) VoteGauge(ctx context.Context, msg *types.MsgVoteGauge) (*types.MsgVoteGaugeResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	totalWeight := math.LegacyZeroDec()
	for _, poolWeight := range msg.PoolWeights {
		weight, err := math.LegacyNewDecFromStr(poolWeight.Weight)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrInvalidWeight, "invalid weight (pool %d): %s", poolWeight.PoolId, err)
		}
		if weight.IsNegative() {
			return nil, errorsmod.Wrapf(types.ErrInvalidWeight, "negative weight (pool %d)", poolWeight.PoolId)
		}
		totalWeight = totalWeight.Add(weight)
	}
	if totalWeight.GT(math.LegacyOneDec()) {
		return nil, errorsmod.Wrapf(types.ErrTotalWeightGTOne, "total weight: %s", totalWeight.String())
	}
	// end static validation

	for _, poolWeight := range msg.PoolWeights {
		_, found, err := k.liquidityPoolKeeper.GetPool(ctx, poolWeight.PoolId)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, liquiditypooltypes.ErrPoolNotFound
		}
	}

	err := k.SetVote(ctx, types.Vote{
		Sender:      msg.Sender,
		PoolWeights: msg.PoolWeights,
	})
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to set vote")
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventVoteGauge{
		Address:     msg.Sender,
		PoolWeights: msg.PoolWeights,
	}); err != nil {
		return nil, err
	}

	return &types.MsgVoteGaugeResponse{}, nil
}
