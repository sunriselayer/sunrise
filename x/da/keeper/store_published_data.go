package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetPublishedData(ctx context.Context, metadataUri string) (data types.PublishedData, found bool, err error) {
	has, err := k.PublishedData.Has(ctx, metadataUri)
	if err != nil {
		return data, false, err
	}

	if !has {
		return data, false, nil
	}

	val, err := k.PublishedData.Get(ctx, metadataUri)
	if err != nil {
		return data, false, err
	}

	return val, true, nil
}

// SetPublishedData sets the published data in the store
func (k Keeper) SetPublishedData(ctx context.Context, data types.PublishedData) error {
	err := k.PublishedData.Set(ctx, data.MetadataUri, data)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeletePublishedData(ctx sdk.Context, data types.PublishedData) error {
	err := k.PublishedData.Remove(ctx, data.MetadataUri)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAllPublishedData(ctx context.Context) (list []types.PublishedData, err error) {
	err = k.PublishedData.Walk(
		ctx,
		nil,
		func(key string, value types.PublishedData) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (k Keeper) GetSpecificStatusData(ctx sdk.Context, status types.Status) (list []types.PublishedData, err error) {
	err = k.PublishedData.Indexes.StatusTime.Walk(
		ctx,
		collections.NewPrefixedPairRange[collections.Pair[string, int64], string](
			collections.PairPrefix[string, int64](status.String()),
		),
		func(key collections.Pair[string, int64], metadataUri string) (bool, error) {
			data, err := k.PublishedData.Get(ctx, metadataUri)
			if err != nil {
				return false, err
			}
			list = append(list, data)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (k Keeper) GetSpecificStatusDataBeforeTime(ctx sdk.Context, status types.Status, timestamp int64) (list []types.PublishedData, err error) {
	err = k.PublishedData.Indexes.StatusTime.Walk(
		ctx,
		collections.NewPrefixedPairRange[collections.Pair[string, int64], string](
			collections.PairPrefix[string, int64](status.String()),
		),
		func(key collections.Pair[string, int64], metadataUri string) (bool, error) {
			if key.K2() > timestamp {
				return true, nil
			}
			data, err := k.PublishedData.Get(ctx, metadataUri)
			if err != nil {
				return false, err
			}
			list = append(list, data)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
