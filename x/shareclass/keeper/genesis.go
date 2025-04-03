package keeper

import (
	"context"
	"time"

	"cosmossdk.io/math"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.Unbondings {
		err := k.SetUnbonding(ctx, elem)
		if err != nil {
			return err
		}
	}
	k.SetUnbondingId(ctx, genState.UnbondingCount)

	for _, elem := range genState.RewardMultipliers {
		validatorAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(elem.Validator)
		if err != nil {
			return err
		}
		rewardMultiplier, err := math.NewDecFromString(elem.RewardMultiplier)
		if err != nil {
			return err
		}

		err = k.SetRewardMultiplier(ctx, validatorAddr, elem.Denom, rewardMultiplier)
		if err != nil {
			return err
		}
	}

	for _, elem := range genState.UserLastRewardMultipliers {
		user, err := k.addressCodec.StringToBytes(elem.User)
		if err != nil {
			return err
		}
		validatorAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(elem.Validator)
		if err != nil {
			return err
		}
		rewardMultiplier, err := math.NewDecFromString(elem.RewardMultiplier)
		if err != nil {
			return err
		}

		err = k.SetUserLastRewardMultiplier(ctx, user, validatorAddr, elem.Denom, rewardMultiplier)
		if err != nil {
			return err
		}
	}

	for _, elem := range genState.LastRewardHandlingTimes {
		validatorAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(elem.Validator)
		if err != nil {
			return err
		}
		lastRewardHandlingTime := time.Unix(elem.LastRewardHandlingTime, 0)

		err = k.SetLastRewardHandlingTime(ctx, validatorAddr, lastRewardHandlingTime)
		if err != nil {
			return err
		}
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
