package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k Keeper) initOrUpdateTick(ctx sdk.Context, poolId uint64, tickIndex int64, liquidityDelta math.LegacyDec, upper bool) (tickIsEmpty bool, err error) {
	tickInfo, err := k.GetTickInfo(ctx, poolId, tickIndex)
	if err != nil {
		return false, err
	}

	liquidityBefore := tickInfo.LiquidityGross
	liquidityAfter := liquidityBefore.Add(liquidityDelta)
	tickInfo.LiquidityGross = liquidityAfter

	if upper {
		tickInfo.LiquidityNet.SubMut(liquidityDelta)
	} else {
		tickInfo.LiquidityNet.AddMut(liquidityDelta)
	}

	if tickInfo.LiquidityGross.IsZero() && tickInfo.LiquidityNet.IsZero() {
		tickIsEmpty = true
	}

	k.SetTickInfo(ctx, tickInfo)
	return tickIsEmpty, nil
}

func (k Keeper) crossTick(ctx sdk.Context, poolId uint64, tickIndex int64, tickInfo *types.TickInfo, swapStateFeeGrowth sdk.DecCoin, feeAccumValue sdk.DecCoins) (err error) {
	if tickInfo == nil {
		return types.ErrNextTickInfoNil
	}

	tickInfo.FeeGrowth = feeAccumValue.
		Add(swapStateFeeGrowth).
		Sub(tickInfo.FeeGrowth)

	// TODO: consider AccumulatorObject

	k.SetTickInfo(ctx, *tickInfo)

	return nil
}

func (k Keeper) newTickInfo(ctx context.Context, poolId uint64, tickIndex int64) (tickStruct types.TickInfo, err error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return tickStruct, types.ErrPoolNotFound
	}
	_ = pool

	// initial fee Growth calculation
	initialFeeGrowth, err := k.getInitialFeeGrowth(ctx, pool, tickIndex)
	if err != nil {
		return tickStruct, err
	}

	// TODO: get initial update growth

	return types.TickInfo{
		PoolId:         poolId,
		TickIndex:      tickIndex,
		LiquidityGross: math.LegacyZeroDec(),
		LiquidityNet:   math.LegacyZeroDec(),
		FeeGrowth:      initialFeeGrowth,
	}, nil
}

// SetTickInfo set a specific tickInfo in the store
func (k Keeper) SetTickInfo(ctx context.Context, tickInfo types.TickInfo) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&tickInfo)
	storeAdapter.Set(types.GetTickInfoIDBytes(tickInfo.PoolId, tickInfo.TickIndex), b)
}

// GetTickInfo returns a tickInfo from its id
func (k Keeper) GetTickInfo(ctx context.Context, poolId uint64, tickIndex int64) (val types.TickInfo, err error) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := storeAdapter.Get(types.GetTickInfoIDBytes(poolId, tickIndex))
	if b == nil {
		return k.newTickInfo(ctx, poolId, tickIndex)
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, nil
}

func (k Keeper) RemoveTickInfo(ctx context.Context, poolId uint64, tickIndex int64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	storeAdapter.Delete(types.GetTickInfoIDBytes(poolId, tickIndex))
}

func (k Keeper) GetAllInitializedTicksForPool(ctx sdk.Context, poolId uint64) (list []types.TickInfo) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(storeAdapter, types.KeyTickPrefixByPoolId(poolId))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TickInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) GetAllTickInfos(ctx context.Context) (list []types.TickInfo) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(storeAdapter, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TickInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func ParseTickFromBz(bz []byte) (tick types.TickInfo, err error) {
	if len(bz) == 0 {
		return types.TickInfo{}, errors.New("tick not found")
	}
	err = proto.Unmarshal(bz, &tick)
	return tick, err
}
