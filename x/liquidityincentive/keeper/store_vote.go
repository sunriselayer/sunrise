package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// SetVote set a specific vote in the store from its index
func (k Keeper) SetVote(ctx context.Context, vote types.Vote) error {
	addr := sdk.MustAccAddressFromBech32(vote.Sender)
	err := k.Votes.Set(ctx, addr, vote)
	if err != nil {
		return err
	}
	return nil
}

// GetVote returns a vote from its index
func (k Keeper) GetVote(ctx context.Context, sender string) (val types.Vote, found bool, err error) {
	addr := sdk.MustAccAddressFromBech32(sender)
	has, err := k.Votes.Has(ctx, addr)
	if err != nil {
		return val, false, err
	}

	if !has {
		return val, false, nil
	}

	val, err = k.Votes.Get(ctx, addr)
	if err != nil {
		return val, false, err
	}

	return val, true, nil
}

// RemoveVote removes a vote from the store
func (k Keeper) RemoveVote(ctx context.Context, sender string) error {
	addr := sdk.MustAccAddressFromBech32(sender)
	err := k.Votes.Remove(ctx, addr)
	if err != nil {
		return err
	}
	return nil
}

// GetAllVotes returns all vote
func (k Keeper) GetAllVotes(ctx context.Context) (list []types.Vote, err error) {
	err = k.Votes.Walk(
		ctx,
		nil,
		func(key sdk.AccAddress, value types.Vote) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (k Keeper) RemoveAllVotes(ctx context.Context) error {
	err := k.Votes.Walk(ctx, nil, func(key sdk.AccAddress, value types.Vote) (bool, error) {
		err := k.Votes.Remove(ctx, key)
		if err != nil {
			return false, err
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}
