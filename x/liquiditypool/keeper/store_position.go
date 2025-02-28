package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// GetPositionCount get the total number of position
func (k Keeper) GetPositionCount(ctx context.Context) (count uint64, err error) {
	count, err = k.PositionId.Peek(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SetPositionCount set the total number of position
func (k Keeper) SetPositionCount(ctx context.Context, count uint64) error {
	err := k.PositionId.Set(ctx, count)
	if err != nil {
		return err
	}

	return nil
}

// AppendPosition appends a position in the store with a new id and update the count
func (k Keeper) AppendPosition(ctx context.Context, position types.Position) (id uint64, err error) {
	// Create the position
	id, err = k.PositionId.Next(ctx)
	if err != nil {
		return 0, err
	}

	// Set the ID of the appended value
	position.Id = id
	err = k.SetPosition(ctx, position)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SetPosition set a specific position in the store
func (k Keeper) SetPosition(ctx context.Context, position types.Position) error {
	err := k.Positions.Set(ctx, position.Id, position)
	if err != nil {
		return err
	}

	return nil
}

// GetPosition returns a position from its id
func (k Keeper) GetPosition(ctx context.Context, id uint64) (val types.Position, found bool, err error) {
	has, err := k.Positions.Has(ctx, id)
	if err != nil {
		return val, false, err
	}

	if !has {
		return val, false, nil
	}

	val, err = k.Positions.Get(ctx, id)
	if err != nil {
		return val, false, err
	}

	return val, true, nil
}

// RemovePosition removes a position from the store
func (k Keeper) RemovePosition(ctx context.Context, id uint64) error {
	err := k.Positions.Remove(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllPositions returns all position
func (k Keeper) GetAllPositions(ctx context.Context) (list []types.Position, err error) {
	err = k.Positions.Walk(
		ctx,
		nil,
		func(key uint64, value types.Position) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (k Keeper) PoolHasPosition(ctx context.Context, poolId uint64) (bool, error) {
	iterator, err := k.Positions.Indexes.PoolId.Iterate(
		ctx,
		collections.NewPrefixedPairRange[uint64, uint64](poolId),
	)
	if err != nil {
		return false, err
	}

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		return true, nil
	}
	return false, nil
}

func (k Keeper) GetPositionsByPool(ctx context.Context, poolId uint64) (list []types.Position, err error) {
	err = k.Positions.Indexes.PoolId.Walk(
		ctx,
		collections.NewPrefixedPairRange[uint64, uint64](poolId),
		func(_ uint64, positionId uint64) (bool, error) {
			value, _, err := k.GetPosition(ctx, positionId)
			if err != nil {
				return false, err
			}
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return
}

func (k Keeper) GetPositionsByAddress(ctx context.Context, addr sdk.AccAddress) (list []types.Position, err error) {
	err = k.Positions.Indexes.Address.Walk(
		ctx,
		collections.NewPrefixedPairRange[sdk.AccAddress, uint64](addr),
		func(_ sdk.AccAddress, positionId uint64) (bool, error) {
			value, _, err := k.GetPosition(ctx, positionId)
			if err != nil {
				return false, err
			}
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return
}
