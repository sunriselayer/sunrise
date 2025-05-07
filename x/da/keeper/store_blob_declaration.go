package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetBlobDeclaration(ctx context.Context, shardsMerkleRoot []byte) (blobDeclaration types.BlobDeclaration, found bool, err error) {

	has, err := k.BlobDeclarations.Has(ctx, shardsMerkleRoot)
	if err != nil {
		return blobDeclaration, false, err
	}

	if !has {
		return blobDeclaration, false, nil
	}

	val, err := k.BlobDeclarations.Get(ctx, shardsMerkleRoot)
	if err != nil {
		return blobDeclaration, false, err
	}

	return val, true, nil
}

// SetBlobDeclaration set the blob declaration of the blob declaration
func (k Keeper) SetBlobDeclaration(ctx context.Context, data types.BlobDeclaration) error {
	err := k.BlobDeclarations.Set(ctx, data.ShardsMerkleRoot, data)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteBlobDeclaration(ctx context.Context, shardsMerkleRoot []byte) error {
	err := k.BlobDeclarations.Remove(ctx, shardsMerkleRoot)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAllBlobDeclarations(ctx context.Context) (list []types.BlobDeclaration, err error) {
	err = k.BlobDeclarations.Walk(
		ctx,
		nil,
		func(key []byte, value types.BlobDeclaration) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
