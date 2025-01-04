package keeper

import (
	"context"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
)

func (k Keeper) InitAccumulator(ctx context.Context, name string) error {
	store := k.KVStoreService.OpenKVStore(ctx)
	hasKey, err := store.Has(types.FormatKeyAccumPrefix(name))
	if err != nil {
		return err
	}
	if hasKey {
		return errors.New("Accumulator with given name already exists in store")
	}

	return k.SetAccumulator(ctx, types.AccumulatorObject{
		Name:        name,
		AccumValue:  sdk.NewDecCoins(),
		TotalShares: math.LegacyZeroDec(),
	})
}

func (k Keeper) GetAccumulator(ctx context.Context, name string) (types.AccumulatorObject, error) {
	accumulator := types.AccumulatorObject{}
	store := k.KVStoreService.OpenKVStore(ctx)
	bz, err := store.Get(types.FormatKeyAccumPrefix(name))
	if err != nil {
		return types.AccumulatorObject{}, err
	}
	if bz == nil {
		return types.AccumulatorObject{}, types.ErrAccumDoesNotExist
	}

	err = proto.Unmarshal(bz, &accumulator)
	if err != nil {
		return types.AccumulatorObject{}, err
	}

	return accumulator, nil
}

func (k Keeper) GetAllAccumulators(ctx context.Context) (list []types.AccumulatorObject) {
	storeAdapter := runtime.KVStoreAdapter(k.KVStoreService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.KeyAccumPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AccumulatorObject
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) SetAccumulator(ctx context.Context, accumulator types.AccumulatorObject) error {
	bz, err := proto.Marshal(&accumulator)
	if err != nil {
		return err
	}
	store := k.KVStoreService.OpenKVStore(ctx)
	return store.Set(types.FormatKeyAccumPrefix(accumulator.Name), bz)
}

func (k Keeper) AddToAccumulator(ctx context.Context, accumulator types.AccumulatorObject, amt sdk.DecCoins) {
	accumulator.AccumValue = accumulator.AccumValue.Add(amt...)
	err := k.SetAccumulator(ctx, accumulator)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetAccumulatorPosition(ctx context.Context, accumName, name string) (types.AccumulatorPosition, error) {
	position := types.AccumulatorPosition{}
	store := k.KVStoreService.OpenKVStore(ctx)
	bz, err := store.Get(types.FormatKeyAccumulatorPositionPrefix(accumName, name))
	if err != nil {
		return types.AccumulatorPosition{}, err
	}
	if bz == nil {
		return types.AccumulatorPosition{}, types.ErrNoPosition
	}

	err = proto.Unmarshal(bz, &position)
	if err != nil {
		return types.AccumulatorPosition{}, err
	}

	return position, nil
}

func (k Keeper) GetAllAccumulatorPositions(ctx context.Context) (list []types.AccumulatorPosition) {
	storeAdapter := runtime.KVStoreAdapter(k.KVStoreService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.KeyAccumulatorPositionPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AccumulatorPosition
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) SetAccumulatorPosition(ctx context.Context, accumName string, accumulatorValuePerShare sdk.DecCoins, index string, numShareUnits math.LegacyDec, unclaimedRewardsTotal sdk.DecCoins) {
	position := types.AccumulatorPosition{
		Name:                  accumName,
		Index:                 index,
		NumShares:             numShareUnits,
		AccumValuePerShare:    accumulatorValuePerShare,
		UnclaimedRewardsTotal: unclaimedRewardsTotal,
	}
	bz, err := proto.Marshal(&position)
	if err != nil {
		panic(err)
	}
	store := k.KVStoreService.OpenKVStore(ctx)
	err = store.Set(types.FormatKeyAccumulatorPositionPrefix(accumName, index), bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) NewPositionIntervalAccumulation(ctx context.Context, accumName, name string, numShareUnits math.LegacyDec, intervalAccumulationPerShare sdk.DecCoins) error {
	k.SetAccumulatorPosition(ctx, accumName, intervalAccumulationPerShare, name, numShareUnits, sdk.NewDecCoins())

	accumulator, err := k.GetAccumulator(ctx, accumName)
	if err != nil {
		return err
	}

	if accumulator.TotalShares.IsNil() {
		accumulator.TotalShares = math.LegacyZeroDec()
	}

	accumulator.TotalShares = accumulator.TotalShares.Add(numShareUnits)
	return k.SetAccumulator(ctx, accumulator)
}

func (k Keeper) AddToPositionIntervalAccumulation(ctx context.Context, accumName, name string, newShares math.LegacyDec, intervalAccumulationPerShare sdk.DecCoins) error {
	if !newShares.IsPositive() {
		return errors.New("Adding non-positive number of shares to position")
	}

	position, err := k.GetAccumulatorPosition(ctx, accumName, name)
	if err != nil {
		return err
	}

	accumulator, err := k.GetAccumulator(ctx, accumName)
	if err != nil {
		return err
	}
	unclaimedRewards := GetTotalRewards(accumulator, position)
	oldNumShares, err := k.GetAccumulatorPositionSize(ctx, accumName, name)
	if err != nil {
		return err
	}

	k.SetAccumulatorPosition(ctx, accumName, intervalAccumulationPerShare, name, oldNumShares.Add(newShares), unclaimedRewards)

	accumulator, err = k.GetAccumulator(ctx, accumName)
	if err != nil {
		return err
	}
	if accumulator.TotalShares.IsNil() {
		accumulator.TotalShares = math.LegacyZeroDec()
	}
	accumulator.TotalShares = accumulator.TotalShares.Add(newShares)
	return k.SetAccumulator(ctx, accumulator)
}

func (k Keeper) RemoveFromPositionIntervalAccumulation(ctx context.Context, accumName, name string, numSharesToRemove math.LegacyDec, intervalAccumulationPerShare sdk.DecCoins) error {
	if !numSharesToRemove.IsPositive() {
		return fmt.Errorf("Removing non-positive shares (%s)", numSharesToRemove)
	}

	position, err := k.GetAccumulatorPosition(ctx, accumName, name)
	if err != nil {
		return err
	}

	if numSharesToRemove.GT(position.NumShares) {
		return fmt.Errorf("Removing more shares (%s) than existing in the position (%s)", numSharesToRemove, position.NumShares)
	}

	accumulator, err := k.GetAccumulator(ctx, accumName)
	if err != nil {
		return err
	}
	unclaimedRewards := GetTotalRewards(accumulator, position)
	oldNumShares, err := k.GetAccumulatorPositionSize(ctx, accumName, name)
	if err != nil {
		return err
	}

	k.SetAccumulatorPosition(ctx, accumName, intervalAccumulationPerShare, name, oldNumShares.Sub(numSharesToRemove), unclaimedRewards)

	accumulator, err = k.GetAccumulator(ctx, accumName)
	if err != nil {
		return err
	}
	if accumulator.TotalShares.IsNil() {
		accumulator.TotalShares = math.LegacyZeroDec()
	}
	accumulator.TotalShares = accumulator.TotalShares.Sub(numSharesToRemove)
	return k.SetAccumulator(ctx, accumulator)
}

func (k Keeper) UpdatePositionIntervalAccumulation(ctx context.Context, accumName, name string, numShares math.LegacyDec, intervalAccumulationPerShare sdk.DecCoins) error {
	if numShares.IsZero() {
		return types.ErrZeroShares
	}

	if numShares.IsNegative() {
		return k.RemoveFromPositionIntervalAccumulation(ctx, accumName, name, numShares.Neg(), intervalAccumulationPerShare)
	}

	return k.AddToPositionIntervalAccumulation(ctx, accumName, name, numShares, intervalAccumulationPerShare)
}

func (k Keeper) SetPositionIntervalAccumulation(ctx context.Context, accumName, name string, intervalAccumulationPerShare sdk.DecCoins) error {
	position, err := k.GetAccumulatorPosition(ctx, accumName, name)
	if err != nil {
		return err
	}

	k.SetAccumulatorPosition(ctx, accumName, intervalAccumulationPerShare, name, position.NumShares, position.UnclaimedRewardsTotal)

	return nil
}

// func (k Keeper) LiquidateAndDeletePosition(ctx context.Context, accumName, positionName string) (sdk.DecCoins, error) {
// 	position, err := k.GetAccumulatorPosition(ctx, accumName, positionName)
// 	if err != nil {
// 		return sdk.DecCoins{}, err
// 	}

// 	remainingRewards, dust, err := k.ClaimRewards(ctx, accumName, positionName)
// 	if err != nil {
// 		return sdk.DecCoins{}, err
// 	}

// 	store := k.storeService.OpenKVStore(ctx)
// 	err = store.Delete(types.FormatKeyAccumulatorPositionPrefix(accumName, positionName))
// 	if err != nil {
// 		return sdk.DecCoins{}, err
// 	}

// 	accumulator, err := k.GetAccumulator(ctx, accumName)
// 	if err != nil {
// 		return sdk.DecCoins{}, err
// 	}
// 	accumulator.TotalShares.SubMut(position.NumShares)
// 	err = k.SetAccumulator(ctx, accumulator)
// 	if err != nil {
// 		return sdk.DecCoins{}, err
// 	}

// 	return sdk.NewDecCoinsFromCoins(remainingRewards...).Add(dust...), nil
// }

func (k Keeper) DeletePosition(ctx context.Context, accumName, positionName string) {
	store := k.KVStoreService.OpenKVStore(ctx)
	err := store.Delete(types.FormatKeyAccumulatorPositionPrefix(accumName, positionName))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetAccumulatorPositionSize(ctx context.Context, accumName, name string) (math.LegacyDec, error) {
	position, err := k.GetAccumulatorPosition(ctx, accumName, name)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return position.NumShares, nil
}

func (k Keeper) HasPosition(ctx context.Context, accumName, name string) bool {
	store := k.KVStoreService.OpenKVStore(ctx)
	containsKey, err := store.Has(types.FormatKeyAccumulatorPositionPrefix(accumName, name))
	if err != nil {
		panic(err)
	}
	return containsKey
}

func (k Keeper) ClaimRewards(ctx context.Context, accumName, positionName string) (sdk.Coins, sdk.DecCoins, error) {
	accumulator, err := k.GetAccumulator(ctx, accumName)
	if err != nil {
		return sdk.Coins{}, sdk.DecCoins{}, types.ErrNoPosition
	}

	position, err := k.GetAccumulatorPosition(ctx, accumName, positionName)
	if err != nil {
		return sdk.Coins{}, sdk.DecCoins{}, types.ErrNoPosition
	}

	totalRewards := GetTotalRewards(accumulator, position)
	truncatedRewardsTotal, dust := totalRewards.TruncateDecimal()

	if position.NumShares.IsZero() {
		k.DeletePosition(ctx, accumName, positionName)
	} else {
		k.SetAccumulatorPosition(ctx, accumName, accumulator.AccumValue, positionName, position.NumShares, sdk.NewDecCoins())
	}

	return truncatedRewardsTotal, dust, nil
}

func GetTotalRewards(accumulator types.AccumulatorObject, position types.AccumulatorPosition) sdk.DecCoins {
	totalRewards := position.UnclaimedRewardsTotal

	if !position.NumShares.IsPositive() {
		return sdk.DecCoins{}
	}
	for _, coin := range position.AccumValuePerShare {
		if accumulator.AccumValue.AmountOf(coin.Denom).LT(coin.Amount) {
			return sdk.DecCoins{}
		}
	}
	accumRewards := accumulator.AccumValue.Sub(position.AccumValuePerShare).MulDec(position.NumShares)
	totalRewards = totalRewards.Add(accumRewards...)

	return totalRewards
}
