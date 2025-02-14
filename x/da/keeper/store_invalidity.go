package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetInvalidity(ctx context.Context, metadataUri string, sender []byte) (invalidity types.Invalidity, found bool) {
	key := collections.Join(metadataUri, sender)
	has, err := k.Invalidities.Has(ctx, key)
	if err != nil {
		panic(err)
	}

	if !has {
		return invalidity, false
	}

	val, err := k.Invalidities.Get(ctx, key)
	if err != nil {
		panic(err)
	}

	return val, true
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

func (k Keeper) DeleteInvalidity(ctx sdk.Context, metadataUri string, sender []byte) {
	err := k.Invalidities.Remove(ctx, collections.Join(metadataUri, sender))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetInvalidities(ctx sdk.Context, metadataUri string) (list []types.Invalidity) {
	err := k.Invalidities.Walk(
		ctx,
		collections.NewPrefixedPairRange[string, []byte](metadataUri),
		func(key collections.Pair[string, []byte], value types.Invalidity) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		panic(err)
	}

	return
}

func (k Keeper) GetAllInvalidities(ctx context.Context) (list []types.Invalidity) {
	err := k.Invalidities.Walk(
		ctx,
		nil,
		func(key collections.Pair[string, []byte], value types.Invalidity) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		panic(err)
	}

	return
}
