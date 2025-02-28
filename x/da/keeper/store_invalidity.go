package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

// GetInvalidity returns an invalidity and whether it was found
func (k Keeper) GetInvalidity(ctx context.Context, metadataUri string, sender []byte) (invalidity types.Invalidity, found bool, err error) {
	key := collections.Join(metadataUri, sender)
	has, err := k.Invalidities.Has(ctx, key)
	if err != nil {
		return invalidity, false, err
	}

	if !has {
		return invalidity, false, nil
	}

	val, err := k.Invalidities.Get(ctx, key)
	if err != nil {
		return invalidity, false, err
	}

	return val, true, nil
}

// SetInvalidity set the Invalidity of the PublishedData
func (k Keeper) SetInvalidity(ctx context.Context, data types.Invalidity) error {
	addr, err := k.addressCodec.StringToBytes(data.Sender)
	if err != nil {
		return err
	}

	err = k.Invalidities.Set(ctx, collections.Join(data.MetadataUri, addr), data)
	if err != nil {
		return err
	}

	return nil
}

// DeleteInvalidity removes an invalidity
func (k Keeper) DeleteInvalidity(ctx sdk.Context, metadataUri string, sender []byte) error {
	return k.Invalidities.Remove(ctx, collections.Join(metadataUri, sender))
}

// GetInvalidities returns all invalidities for a metadata URI
func (k Keeper) GetInvalidities(ctx sdk.Context, metadataUri string) (list []types.Invalidity, err error) {
	err = k.Invalidities.Walk(
		ctx,
		collections.NewPrefixedPairRange[string, []byte](metadataUri),
		func(key collections.Pair[string, []byte], value types.Invalidity) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetAllInvalidities returns all invalidities
func (k Keeper) GetAllInvalidities(ctx context.Context) (list []types.Invalidity, err error) {
	err = k.Invalidities.Walk(
		ctx,
		nil,
		func(key collections.Pair[string, []byte], value types.Invalidity) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return list, nil
}
