package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetProof(ctx context.Context, metadataUri string, sender []byte) (data types.Proof) {
	key := collections.Join(metadataUri, sender)
	has, err := k.Proofs.Has(ctx, key)
	if err != nil {
		panic(err)
	}

	if !has {
		return data
	}

	val, err := k.Proofs.Get(ctx, key)
	if err != nil {
		panic(err)
	}

	return val
}

// SetProof set the proof of the PublishedData
func (k Keeper) SetProof(ctx context.Context, data types.Proof) error {
	addr, err := k.addressCodec.StringToBytes(data.Sender)
	if err != nil {
		return err
	}

	err = k.Proofs.Set(ctx, collections.Join(data.MetadataUri, addr), data)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteProof(ctx sdk.Context, metadataUri string, sender []byte) {
	err := k.Proofs.Remove(ctx, collections.Join(metadataUri, sender))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetProofs(ctx sdk.Context, metadataUri string) (list []types.Proof) {
	err := k.Proofs.Walk(
		ctx,
		collections.NewPrefixedPairRange[string, []byte](metadataUri),
		func(key collections.Pair[string, []byte], value types.Proof) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		panic(err)
	}

	return
}

func (k Keeper) GetAllProofs(ctx context.Context) (list []types.Proof) {
	err := k.Proofs.Walk(
		ctx,
		nil,
		func(key collections.Pair[string, []byte], value types.Proof) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		panic(err)
	}

	return
}
