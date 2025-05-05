package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func keyBlobDeclaration(blockHeight int64, shardsMerkleRoot []byte) collections.Pair[int64, []byte] {
	return collections.Join(blockHeight, shardsMerkleRoot)
}

func (k Keeper) GetBlobDeclaration(ctx context.Context, blockHeight int64, shardsMerkleRoot []byte) (blobDeclaration types.BlobDeclaration, found bool, err error) {
	key := keyBlobDeclaration(blockHeight, shardsMerkleRoot)
	has, err := k.BlobDeclarations.Has(ctx, key)
	if err != nil {
		return blobDeclaration, false, err
	}

	if !has {
		return blobDeclaration, false, nil
	}

	val, err := k.BlobDeclarations.Get(ctx, key)
	if err != nil {
		return blobDeclaration, false, err
	}

	return val, true, nil
}

// SetBlobDeclaration set the blob declaration of the blob declaration
func (k Keeper) SetBlobDeclaration(ctx context.Context, data types.BlobDeclaration) error {
	err := k.BlobDeclarations.Set(ctx, keyBlobDeclaration(data.BlockHeight, data.ShardsMerkleRoot), data)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteBlobDeclaration(ctx context.Context, declarationHeight int64, shardsMerkleRoot []byte) error {
	err := k.BlobDeclarations.Remove(ctx, keyBlobDeclaration(declarationHeight, shardsMerkleRoot))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAllBlobDeclarations(ctx context.Context) (list []types.BlobDeclaration, err error) {
	err = k.BlobDeclarations.Walk(
		ctx,
		nil,
		func(key collections.Pair[int64, []byte], value types.BlobDeclaration) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
