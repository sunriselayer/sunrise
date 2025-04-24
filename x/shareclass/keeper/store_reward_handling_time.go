package keeper

import (
	"context"
	"errors"
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLastRewardHandlingTime(ctx context.Context, validatorAddr sdk.ValAddress) (time.Time, error) {
	lastRewardHandlingTime, err := k.LastRewardHandlingTime.Get(ctx, validatorAddr)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// Return default time (Unix time 0)
			return time.Unix(0, 0), nil
		}
		return time.Time{}, err
	}

	return time.Unix(lastRewardHandlingTime, 0), nil
}

func (k Keeper) SetLastRewardHandlingTime(ctx context.Context, validatorAddr sdk.ValAddress, lastRewardHandlingTime time.Time) error {
	return k.LastRewardHandlingTime.Set(ctx, validatorAddr, lastRewardHandlingTime.Unix())
}
