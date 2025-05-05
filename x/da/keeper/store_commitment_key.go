package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetCommitmentKey(ctx context.Context, validator sdk.ValAddress) (commitmentKey types.CommitmentKey, found bool, err error) {
	has, err := k.CommitmentKeys.Has(ctx, validator)
	if err != nil {
		return commitmentKey, false, err
	}

	if !has {
		return commitmentKey, false, nil
	}

	commitmentKey, err = k.CommitmentKeys.Get(ctx, validator)
	if err != nil {
		return commitmentKey, false, err
	}

	return commitmentKey, true, nil
}

// SetCommitmentKey set the commitment key of the validator
func (k Keeper) SetCommitmentKey(ctx context.Context, validator sdk.ValAddress, commitmentKey types.CommitmentKey) error {
	err := k.CommitmentKeys.Set(ctx, validator, commitmentKey)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteCommitmentKey(ctx context.Context, validator sdk.ValAddress) error {
	err := k.CommitmentKeys.Remove(ctx, validator)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAllCommitmentKeys(ctx context.Context) (list []types.CommitmentKey, err error) {
	err = k.CommitmentKeys.Walk(
		ctx,
		nil,
		func(key []byte, value types.CommitmentKey) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
