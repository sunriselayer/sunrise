package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SetTwap set a specific twap in the store from its index
func (k Keeper) SetTwap(ctx context.Context, twap types.Twap) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TwapKeyPrefix))
	b := k.cdc.MustMarshal(&twap)
	store.Set(types.TwapKey(
		twap.BaseDenom,
		twap.QuoteDenom,
	), b)
}

// GetTwap returns a twap from its index
func (k Keeper) GetTwap(
	ctx context.Context,
	baseDenom string,
	quoteDenom string,
) (val types.Twap, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TwapKeyPrefix))

	b := store.Get(types.TwapKey(
		baseDenom,
		quoteDenom,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveTwap removes a twap from the store
func (k Keeper) RemoveTwap(
	ctx context.Context,
	baseDenom string,
	quoteDenom string,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TwapKeyPrefix))
	store.Delete(types.TwapKey(
		baseDenom,
		quoteDenom,
	))
}

// GetAllTwap returns all twap
func (k Keeper) GetAllTwap(ctx context.Context) (list []types.Twap) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TwapKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Twap
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) RecordTrade(ctx context.Context, baseDenom string, quoteDenom string, price math.LegacyDec, volume math.Int) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if volume.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "volume cannot be zero")
	}

	footprint := types.TradeFootprint{
		Price:     price,
		Volume:    volume,
		Timestamp: sdkCtx.BlockTime(),
	}
	k.SetTradeFootprint(ctx, baseDenom, quoteDenom, footprint)

	footprints := k.GetAllTradeFootprint(ctx, baseDenom, quoteDenom)
	filteredFootprints := []types.TradeFootprint{}

	for _, footprint := range footprints {
		if footprint.Timestamp.After(sdkCtx.BlockTime().Add(-k.GetParams(ctx).TwapWindow)) {
			filteredFootprints = append(filteredFootprints, footprint)
		} else {
			k.RemoveTradeFootprint(ctx, baseDenom, quoteDenom, footprint.Timestamp)
		}
	}

	twap, err := types.CalculateTwap(filteredFootprints)
	if err != nil {
		return err
	}

	timestamp := sdkCtx.BlockTime()
	k.SetTwap(ctx, types.Twap{
		BaseDenom:  baseDenom,
		QuoteDenom: quoteDenom,
		Value:      twap,
		Timestamp:  &timestamp,
	})

	return nil
}

func (k Keeper) GetValidTwap(ctx context.Context, baseDenom string, quoteDenom string) *math.LegacyDec {
	twap, found := k.GetTwap(ctx, baseDenom, quoteDenom)
	if !found {
		return nil
	}

	if twap.Value == nil || twap.Timestamp == nil {
		return nil

	}

	if twap.Timestamp.Before(sdk.UnwrapSDKContext(ctx).BlockTime().Add(-k.GetParams(ctx).TwapExpiry)) {
		return nil
	}

	return twap.Value
}

func (k Keeper) GetValidTwapByPoolId(ctx context.Context, poolId uint64) (*types.Pool, *math.LegacyDec) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, nil
	}

	return &pool, k.GetValidTwap(ctx, pool.BaseDenom, pool.QuoteDenom)
}

func (k Keeper) GetValidSrValueByPoolId(ctx context.Context, poolId uint64, srDenom string) *math.LegacyDec {
	pool, twap := k.GetValidTwapByPoolId(ctx, poolId)
	if pool == nil || twap == nil {
		return nil
	}

	twap2 := k.GetValidTwap(ctx, pool.QuoteDenom, srDenom)
	if twap2 == nil {
		return nil
	}

	result := twap.Mul(*twap2)

	return &result
}
