package keeper

import (
	"bytes"
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/da/types"
)

// GetChallengeCount get the total number of pool
func (k Keeper) GetChallengeCount(ctx context.Context) (count uint64, err error) {
	count, err = k.ChallengeId.Peek(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SetChallengeCount set the total number of pool
func (k Keeper) SetChallengeCount(ctx context.Context, count uint64) error {
	err := k.ChallengeId.Set(ctx, count)
	if err != nil {
		return err
	}
	return nil
}

// AppendChallenge appends a pool in the store with a new id and update the count
func (k Keeper) AppendChallenge(ctx context.Context, pool types.Challenge) (id uint64, err error) {
	// Create the pool
	id, err = k.ChallengeId.Next(ctx)
	if err != nil {
		return 0, err
	}

	// Set the ID of the appended value
	pool.Id = id
	err = k.SetChallenge(ctx, pool)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (k Keeper) GetChallenge(ctx context.Context, id uint64) (data types.Challenge, found bool, err error) {
	has, err := k.Challenges.Has(ctx, id)
	if err != nil {
		return data, false, err
	}

	if !has {
		return data, false, nil
	}

	val, err := k.Challenges.Get(ctx, id)
	if err != nil {
		return data, false, err
	}

	return val, true, nil
}

func (k Keeper) SetChallenge(ctx context.Context, data types.Challenge) error {
	err := k.Challenges.Set(ctx, data.Id, data)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteChallenge(ctx context.Context, id uint64) error {
	err := k.Challenges.Remove(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAllChallenges(ctx context.Context) (list []types.Challenge, err error) {
	err = k.Challenges.Walk(
		ctx,
		nil,
		func(key uint64, value types.Challenge) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (k Keeper) GetAllChallengesByShardsMerkleRoot(ctx context.Context, shardsMerkleRoot []byte) (list []types.Challenge, err error) {
	err = k.Challenges.Indexes.ShardsMerkleRoot.Walk(
		ctx,
		nil,
		func(indexingKey collections.Pair[[]byte, uint64], indexedKey uint64) (stop bool, err error) {
			// TODO: use range prefix
			if !bytes.Equal(indexingKey.K1(), shardsMerkleRoot) {
				return false, nil
			}

			val, err := k.Challenges.Get(ctx, indexedKey)
			if err != nil {
				return true, err
			}

			list = append(list, val)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
