package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) addDenomFromCreator(ctx sdk.Context, creator, denom string) {
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		panic(err)
	}
	key := collections.Join(creatorAddr, denom)
	_ = k.DenomFromCreator.Set(ctx, key, []byte{})
	_ = k.CreatorAddresses.Set(ctx, denom, creatorAddr)
}

func (k Keeper) getDenomsFromCreator(ctx sdk.Context, creator string) []string {
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		panic(err)
	}

	var denoms []string
	err = k.DenomFromCreator.Walk(ctx, collections.NewPrefixedPairRange[sdk.AccAddress, string](creatorAddr), func(key collections.Pair[sdk.AccAddress, string], value []byte) (bool, error) {
		denoms = append(denoms, key.K2())
		return false, nil
	})

	if err != nil {
		panic(err)
	}

	return denoms
}

func (k Keeper) GetAllDenomsIterator(ctx context.Context) (collections.Iterator[collections.Pair[sdk.AccAddress, string], []byte], error) {
	return k.DenomFromCreator.Iterate(ctx, nil)
}
