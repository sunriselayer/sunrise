package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

// GetProof returns a proof and whether it was found for the given metadata URI and sender
func (k Keeper) GetProof(ctx context.Context, metadataUri string, sender []byte) (proof types.Proof, found bool, err error) {
	key := collections.Join(metadataUri, sender)
	has, err := k.Proofs.Has(ctx, key)
	if err != nil {
		return proof, false, err
	}

	if !has {
		return proof, false, nil
	}

	val, err := k.Proofs.Get(ctx, key)
	if err != nil {
		return proof, false, err
	}

	return val, true, nil
}

// SetProof sets the proof of the PublishedData
func (k Keeper) SetProof(ctx context.Context, data types.Proof) error {
	addr, err := k.addressCodec.StringToBytes(data.Sender)
	if err != nil {
		return err
	}

	return k.Proofs.Set(ctx, collections.Join(data.MetadataUri, addr), data)
}

// DeleteProof removes a proof for the given metadata URI and sender
func (k Keeper) DeleteProof(ctx sdk.Context, metadataUri string, sender []byte) error {
	return k.Proofs.Remove(ctx, collections.Join(metadataUri, sender))
}

// GetProofs returns all proofs for a given metadata URI
func (k Keeper) GetProofs(ctx sdk.Context, metadataUri string) (list []types.Proof, err error) {
	err = k.Proofs.Walk(
		ctx,
		collections.NewPrefixedPairRange[string, []byte](metadataUri),
		func(key collections.Pair[string, []byte], value types.Proof) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetAllProofs returns all proofs in the store
func (k Keeper) GetAllProofs(ctx context.Context) (list []types.Proof, err error) {
	err = k.Proofs.Walk(
		ctx,
		nil,
		func(key collections.Pair[string, []byte], value types.Proof) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
