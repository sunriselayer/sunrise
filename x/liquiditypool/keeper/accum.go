package keeper

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
)

type AccumulatorObject struct {
	store         store.KVStore
	name          string
	valuePerShare sdk.DecCoins
	totalShares   math.LegacyDec
}

func MakeAccumulator(accumStore store.KVStore, accumName string) error {
	hasKey, err := accumStore.Has(types.FormatKeyAccumPrefix(accumName))
	if err != nil {
		return err
	}
	if hasKey {
		return errors.New("Accumulator with given name already exists in store")
	}

	initAccumValue := sdk.NewDecCoins()
	initTotalShares := math.LegacyZeroDec()

	newAccum := &AccumulatorObject{accumStore, accumName, initAccumValue, initTotalShares}

	return setAccumulator(newAccum, initAccumValue, initTotalShares)
}

func MakeAccumulatorWithValueAndShare(accumStore store.KVStore, accumName string, accumValue sdk.DecCoins, totalShares math.LegacyDec) error {
	hasKey, err := accumStore.Has(types.FormatKeyAccumPrefix(accumName))
	if err != nil {
		return err
	}
	if hasKey {
		return errors.New("Accumulator with given name already exists in store")
	}

	newAccum := AccumulatorObject{accumStore, accumName, accumValue, totalShares}

	return setAccumulator(&newAccum, accumValue, totalShares)
}

func OverwriteAccumulatorUnsafe(accumStore store.KVStore, accumName string, accumValue sdk.DecCoins, totalShares math.LegacyDec) error {
	hasKey, err := accumStore.Has(types.FormatKeyAccumPrefix(accumName))
	if err != nil {
		return err
	}
	if !hasKey {
		return errors.New("Accumulator with given name does not exist in store")
	}

	newAccum := AccumulatorObject{accumStore, accumName, accumValue, totalShares}

	return setAccumulator(&newAccum, accumValue, totalShares)
}

func GetAccumulator(accumStore store.KVStore, accumName string) (*AccumulatorObject, error) {
	accumContent := types.AccumulatorContent{}
	bz, err := accumStore.Get(types.FormatKeyAccumPrefix(accumName))
	if err != nil {
		return &AccumulatorObject{}, err
	}
	if bz == nil {
		return &AccumulatorObject{}, types.ErrAccumDoesNotExist
	}

	accum := AccumulatorObject{accumStore, accumName, accumContent.AccumValue, accumContent.TotalShares}

	return &accum, nil
}

func (accum AccumulatorObject) MustGetPosition(name string) types.Record {
	position := types.Record{}
	bz, err := accum.store.Get(types.FormatKeyAccumPositionPrefix(accum.name, name))
	if err != nil {
		panic(err)
	}
	if bz == nil {
		panic("empty position")
	}
	err = proto.Unmarshal(bz, &position)
	if err != nil {
		panic(err)
	}
	return position
}

func (accum AccumulatorObject) GetPosition(name string) (types.Record, error) {
	position := types.Record{}
	bz, err := accum.store.Get(types.FormatKeyAccumPositionPrefix(accum.name, name))
	if err != nil {
		return position, err
	}
	if bz == nil {
		panic("empty position")
	}
	err = proto.Unmarshal(bz, &position)
	if err != nil {
		return position, err
	}
	return position, nil
}

func setAccumulator(accum *AccumulatorObject, value sdk.DecCoins, shares math.LegacyDec) error {
	newAccum := types.AccumulatorContent{
		AccumValue:  value,
		TotalShares: shares,
	}
	bz, err := proto.Marshal(&newAccum)
	if err != nil {
		panic(err)
	}
	return accum.store.Set(types.FormatKeyAccumPrefix(accum.name), bz)
}

func (accum *AccumulatorObject) AddToAccumulator(amt sdk.DecCoins) {
	accum.valuePerShare = accum.valuePerShare.Add(amt...)
	err := setAccumulator(accum, accum.valuePerShare, accum.totalShares)
	if err != nil {
		panic(err)
	}
}

func (accum *AccumulatorObject) NewPosition(name string, numShareUnits math.LegacyDec) error {
	return accum.NewPositionIntervalAccumulation(name, numShareUnits, accum.valuePerShare)
}

func (accum *AccumulatorObject) NewPositionIntervalAccumulation(name string, numShareUnits math.LegacyDec, intervalAccumulationPerShare sdk.DecCoins) error {
	initOrUpdatePosition(accum, intervalAccumulationPerShare, name, numShareUnits, sdk.NewDecCoins())

	updatedAccum, err := GetAccumulator(accum.store, accum.name)
	if err != nil {
		return err
	}

	if updatedAccum.totalShares.IsNil() {
		updatedAccum.totalShares = math.LegacyZeroDec()
	}

	accum.totalShares = updatedAccum.totalShares.Add(numShareUnits)
	return setAccumulator(accum, accum.valuePerShare, accum.totalShares)
}

func (accum *AccumulatorObject) AddToPosition(name string, newShares math.LegacyDec) error {
	return accum.AddToPositionIntervalAccumulation(name, newShares, accum.valuePerShare)
}

func (accum *AccumulatorObject) AddToPositionIntervalAccumulation(name string, newShares math.LegacyDec, intervalAccumulationPerShare sdk.DecCoins) error {
	if !newShares.IsPositive() {
		return errors.New("Adding non-positive number of shares to position")
	}

	position, err := GetPosition(accum, name)
	if err != nil {
		return err
	}

	unclaimedRewards := GetTotalRewards(accum, position)
	oldNumShares, err := accum.GetPositionSize(name)
	if err != nil {
		return err
	}

	initOrUpdatePosition(accum, intervalAccumulationPerShare, name, oldNumShares.Add(newShares), unclaimedRewards)

	updatedAccum, err := GetAccumulator(accum.store, accum.name)
	if err != nil {
		return err
	}
	accum.totalShares = updatedAccum.totalShares.Add(newShares)
	return setAccumulator(accum, accum.valuePerShare, accum.totalShares)
}

func (accum *AccumulatorObject) RemoveFromPosition(name string, numSharesToRemove math.LegacyDec) error {
	return accum.RemoveFromPositionIntervalAccumulation(name, numSharesToRemove, accum.valuePerShare)
}

func (accum *AccumulatorObject) RemoveFromPositionIntervalAccumulation(name string, numSharesToRemove math.LegacyDec, intervalAccumulationPerShare sdk.DecCoins) error {
	if !numSharesToRemove.IsPositive() {
		return fmt.Errorf("Removing non-positive shares (%s)", numSharesToRemove)
	}

	position, err := GetPosition(accum, name)
	if err != nil {
		return err
	}

	if numSharesToRemove.GT(position.NumShares) {
		return fmt.Errorf("Removing more shares (%s) than existing in the position (%s)", numSharesToRemove, position.NumShares)
	}

	unclaimedRewards := GetTotalRewards(accum, position)
	oldNumShares, err := accum.GetPositionSize(name)
	if err != nil {
		return err
	}

	initOrUpdatePosition(accum, intervalAccumulationPerShare, name, oldNumShares.Sub(numSharesToRemove), unclaimedRewards)

	updatedAccum, err := GetAccumulator(accum.store, accum.name)
	if err != nil {
		return err
	}
	accum.totalShares = updatedAccum.totalShares.Sub(numSharesToRemove)
	return setAccumulator(accum, accum.valuePerShare, accum.totalShares)
}
func (accum *AccumulatorObject) UpdatePosition(name string, numShares math.LegacyDec) error {
	return accum.UpdatePositionIntervalAccumulation(name, numShares, accum.valuePerShare)
}

func (accum *AccumulatorObject) UpdatePositionIntervalAccumulation(name string, numShares math.LegacyDec, intervalAccumulationPerShare sdk.DecCoins) error {
	if numShares.IsZero() {
		return types.ErrZeroShares
	}

	if numShares.IsNegative() {
		return accum.RemoveFromPositionIntervalAccumulation(name, numShares.Neg(), intervalAccumulationPerShare)
	}

	return accum.AddToPositionIntervalAccumulation(name, numShares, intervalAccumulationPerShare)
}

func (accum *AccumulatorObject) SetPositionIntervalAccumulation(name string, intervalAccumulationPerShare sdk.DecCoins) error {
	position, err := GetPosition(accum, name)
	if err != nil {
		return err
	}

	initOrUpdatePosition(accum, intervalAccumulationPerShare, name, position.NumShares, position.UnclaimedRewardsTotal)

	return nil
}

func (accum *AccumulatorObject) DeletePosition(positionName string) (sdk.DecCoins, error) {
	position, err := accum.GetPosition(positionName)
	if err != nil {
		return sdk.DecCoins{}, err
	}

	remainingRewards, dust, err := accum.ClaimRewards(positionName)
	if err != nil {
		return sdk.DecCoins{}, err
	}

	err = accum.store.Delete(types.FormatKeyAccumPositionPrefix(accum.name, positionName))
	if err != nil {
		return sdk.DecCoins{}, err
	}

	accum.totalShares.SubMut(position.NumShares)
	err = setAccumulator(accum, accum.valuePerShare, accum.totalShares)
	if err != nil {
		return sdk.DecCoins{}, err
	}

	return sdk.NewDecCoinsFromCoins(remainingRewards...).Add(dust...), nil
}

func (accum AccumulatorObject) deletePosition(positionName string) {
	err := accum.store.Delete(types.FormatKeyAccumPositionPrefix(accum.name, positionName))
	if err != nil {
		panic(err)
	}
}

func (accum *AccumulatorObject) GetPositionSize(name string) (math.LegacyDec, error) {
	position, err := GetPosition(accum, name)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return position.NumShares, nil
}

func (accum AccumulatorObject) HasPosition(name string) bool {
	containsKey, err := accum.store.Has(types.FormatKeyAccumPositionPrefix(accum.name, name))
	if err != nil {
		panic(err)
	}
	return containsKey
}

func (accum AccumulatorObject) GetName() string {
	return accum.name
}

func (accum AccumulatorObject) GetValue() sdk.DecCoins {
	return accum.valuePerShare
}

func (accum *AccumulatorObject) ClaimRewards(positionName string) (sdk.Coins, sdk.DecCoins, error) {
	position, err := GetPosition(accum, positionName)
	if err != nil {
		return sdk.Coins{}, sdk.DecCoins{}, types.ErrNoPosition
	}

	totalRewards := GetTotalRewards(accum, position)
	truncatedRewardsTotal, dust := totalRewards.TruncateDecimal()

	if position.NumShares.IsZero() {
		accum.deletePosition(positionName)
	} else {
		initOrUpdatePosition(accum, accum.valuePerShare, positionName, position.NumShares, sdk.NewDecCoins())
	}

	return truncatedRewardsTotal, dust, nil
}

func (accum AccumulatorObject) GetTotalShares() math.LegacyDec {
	return accum.totalShares
}

func (accum *AccumulatorObject) AddToUnclaimedRewards(positionName string, rewardsToAddTotal sdk.DecCoins) error {
	position, err := GetPosition(accum, positionName)
	if err != nil {
		return err
	}

	if rewardsToAddTotal.IsAnyNegative() {
		return types.ErrNegRewardAddition
	}

	initOrUpdatePosition(accum, position.AccumValuePerShare, positionName, position.NumShares, position.UnclaimedRewardsTotal.Add(rewardsToAddTotal...))

	return nil
}

func initOrUpdatePosition(accum *AccumulatorObject, accumulatorValuePerShare sdk.DecCoins, index string, numShareUnits math.LegacyDec, unclaimedRewardsTotal sdk.DecCoins) {
	position := types.Record{
		NumShares:             numShareUnits,
		AccumValuePerShare:    accumulatorValuePerShare,
		UnclaimedRewardsTotal: unclaimedRewardsTotal,
	}
	bz, err := proto.Marshal(&position)
	if err != nil {
		panic(err)
	}
	err = accum.store.Set(types.FormatKeyAccumPositionPrefix(accum.name, index), bz)
	if err != nil {
		panic(err)
	}
}

func GetPosition(accum *AccumulatorObject, name string) (types.Record, error) {
	position := types.Record{}
	bz, err := accum.store.Get(types.FormatKeyAccumPositionPrefix(accum.name, name))
	if err != nil {
		return types.Record{}, err
	}
	if bz == nil {
		return types.Record{}, types.ErrNoPosition
	}

	err = proto.Unmarshal(bz, &position)
	if err != nil {
		return types.Record{}, err
	}

	return position, nil
}

func GetTotalRewards(accum *AccumulatorObject, position types.Record) sdk.DecCoins {
	totalRewards := position.UnclaimedRewardsTotal

	accumulatorRewards := accum.valuePerShare.Sub(position.AccumValuePerShare).MulDec(position.NumShares)
	totalRewards = totalRewards.Add(accumulatorRewards...)

	return totalRewards
}
