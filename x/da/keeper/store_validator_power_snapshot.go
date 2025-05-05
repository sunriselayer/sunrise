package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func keyValidatorPowerSnapshot(blockHeight int64, validatorAddr []byte) collections.Pair[int64, []byte] {
	return collections.Join(blockHeight, validatorAddr)
}

func (k Keeper) GetValidatorPowerSnapshot(ctx context.Context, blockHeight int64, validator sdk.ValAddress) (data types.ValidatorPowerSnapshot, found bool, err error) {
	has, err := k.ValidatorPowerSnapshots.Has(ctx, keyValidatorPowerSnapshot(blockHeight, validator))
	if err != nil {
		return data, false, err
	}

	if !has {
		return data, false, nil
	}

	val, err := k.ValidatorPowerSnapshots.Get(ctx, keyValidatorPowerSnapshot(blockHeight, validator))
	if err != nil {
		return data, false, err
	}

	return val, true, nil
}

// SetParams set the params
func (k Keeper) SetValidatorPowerSnapshot(ctx context.Context, data types.ValidatorPowerSnapshot) error {
	validatorAddr, err := k.StakingKeeper.ValidatorAddressCodec().StringToBytes(data.Validator)
	if err != nil {
		return err
	}
	err = k.ValidatorPowerSnapshots.Set(ctx, keyValidatorPowerSnapshot(data.BlockHeight, validatorAddr), data)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteValidatorPowerSnapshot(ctx context.Context, blockHeight int64, validator sdk.ValAddress) error {
	err := k.ValidatorPowerSnapshots.Remove(ctx, keyValidatorPowerSnapshot(blockHeight, validator))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAllValidatorPowerSnapshots(ctx context.Context) (list []types.ValidatorPowerSnapshot, err error) {
	err = k.ValidatorPowerSnapshots.Walk(
		ctx,
		nil,
		func(key collections.Pair[int64, []byte], value types.ValidatorPowerSnapshot) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (k Keeper) GetAllValidatorPowerSnapshotsByBlockHeight(ctx context.Context, blockHeight int64) (list []types.ValidatorPowerSnapshot, err error) {
	err = k.ValidatorPowerSnapshots.Walk(
		ctx,
		collections.NewPrefixedPairRange[int64, []byte](blockHeight),
		func(key collections.Pair[int64, []byte], value types.ValidatorPowerSnapshot) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
