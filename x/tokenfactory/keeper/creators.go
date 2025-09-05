package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) addDenomFromCreator(ctx sdk.Context, creator, denom string) error {
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return err
	}
	key := collections.Join(creatorAddr, denom)
	err = k.DenomFromCreator.Set(ctx, key, []byte{})
	if err != nil {
		return err
	}
	err = k.CreatorAddresses.Set(ctx, denom, creatorAddr)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) getDenomsFromCreator(ctx sdk.Context, creator string) ([]string, error) {
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return nil, err
	}

	var denoms []string
	err = k.DenomFromCreator.Walk(ctx, collections.NewPrefixedPairRange[sdk.AccAddress, string](creatorAddr), func(key collections.Pair[sdk.AccAddress, string], value []byte) (bool, error) {
		denoms = append(denoms, key.K2())
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	return denoms, nil
}

func (k Keeper) GetAllDenomsIterator(ctx context.Context) (collections.Iterator[collections.Pair[sdk.AccAddress, string], []byte], error) {
	return k.DenomFromCreator.Iterate(ctx, nil)
}
