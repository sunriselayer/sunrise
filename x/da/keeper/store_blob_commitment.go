package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetBlobCommitment(ctx context.Context, shardsMerkleRoot []byte) (blobCommitment types.BlobCommitment, found bool, err error) {
	has, err := k.BlobCommitments.Has(ctx, shardsMerkleRoot)
	if err != nil {
		return blobCommitment, false, err
	}

	if !has {
		return blobCommitment, false, nil
	}

	val, err := k.BlobCommitments.Get(ctx, shardsMerkleRoot)
	if err != nil {
		return blobCommitment, false, err
	}

	return val, true, nil
}

// SetBlobCommitment set the BlobCommitment of the BlobCommitment
func (k Keeper) SetBlobCommitment(ctx context.Context, data types.BlobCommitment) error {
	err := k.BlobCommitments.Set(ctx, data.ShardsMerkleRoot, data)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteBlobCommitment(ctx context.Context, shardsMerkleRoot []byte) error {
	err := k.BlobCommitments.Remove(ctx, shardsMerkleRoot)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAllBlobCommitments(ctx context.Context) (list []types.BlobCommitment, err error) {
	err = k.BlobCommitments.Walk(
		ctx,
		nil,
		func(key []byte, value types.BlobCommitment) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
