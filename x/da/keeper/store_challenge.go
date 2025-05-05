package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func keyChallenge(shardsMerkleRoot []byte, shardIndex uint32, evaluationPoint uint32) collections.Triple[[]byte, uint32, uint32] {
	return collections.Join3(shardsMerkleRoot, shardIndex, evaluationPoint)
}

func (k Keeper) GetChallenge(ctx context.Context, shardsMerkleRoot []byte, shardIndex uint32, evaluationPoint uint32) (data types.Challenge, found bool, err error) {
	has, err := k.Challenges.Has(ctx, keyChallenge(shardsMerkleRoot, shardIndex, evaluationPoint))
	if err != nil {
		return data, false, err
	}

	if !has {
		return data, false, nil
	}

	val, err := k.Challenges.Get(ctx, keyChallenge(shardsMerkleRoot, shardIndex, evaluationPoint))
	if err != nil {
		return data, false, err
	}

	return val, true, nil
}

func (k Keeper) SetChallenge(ctx context.Context, data types.Challenge) error {
	err := k.Challenges.Set(ctx, keyChallenge(data.ShardsMerkleRoot, data.ShardIndex, data.EvaluationPointIndex), data)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteChallenge(ctx context.Context, shardsMerkleRoot []byte, shardIndex uint32, evaluationPointIndex uint32) error {
	err := k.Challenges.Remove(ctx, keyChallenge(shardsMerkleRoot, shardIndex, evaluationPointIndex))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAllChallenges(ctx context.Context) (list []types.Challenge, err error) {
	err = k.Challenges.Walk(
		ctx,
		nil,
		func(key collections.Triple[[]byte, uint32, uint32], value types.Challenge) (bool, error) {
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
	err = k.Challenges.Walk(
		ctx,
		collections.NewPrefixedTripleRange[[]byte, uint32, uint32](shardsMerkleRoot),
		func(key collections.Triple[[]byte, uint32, uint32], value types.Challenge) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
