package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
	err := k.SetUnbondingId(ctx, genState.UnbondingCount)
	if err != nil {
		return err
	}

	for _, elem := range genState.RewardMultipliers {
		validatorAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(elem.Validator)
		if err != nil {
			return err
		}
		rewardMultiplier, err := math.LegacyNewDecFromStr(elem.RewardMultiplier)
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
		rewardMultiplier, err := math.LegacyNewDecFromStr(elem.RewardMultiplier)
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
// Unbondings and reward indexes are included so a shutdown snapshot can rebuild holder records.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	genesis := types.DefaultGenesis()
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("ExportGenesis: get params: %w", err)
	}
	genesis.Params = params

	genesis.Unbondings, err = k.GetAllUnbondings(ctx)
	if err != nil {
		return nil, fmt.Errorf("ExportGenesis: get unbondings: %w", err)
	}
	genesis.UnbondingCount, err = k.GetUnbondingId(ctx)
	if err != nil {
		return nil, fmt.Errorf("ExportGenesis: get unbonding count: %w", err)
	}

	if err := k.exportRewardState(ctx, genesis); err != nil {
		return nil, err
	}
	return genesis, nil
}

// exportRewardState copies reward indexes into genesis.
// Walk order follows the store key order so the JSON is the same on every node.
func (k Keeper) exportRewardState(ctx context.Context, genesis *types.GenesisState) error {
	err := k.RewardMultiplier.Walk(ctx, nil, func(key collections.Pair[[]byte, string], value string) (bool, error) {
		validator, err := k.stakingKeeper.ValidatorAddressCodec().BytesToString(key.K1())
		if err != nil {
			return true, fmt.Errorf("ExportGenesis: validator address for reward multiplier denom %s: %w", key.K2(), err)
		}
		genesis.RewardMultipliers = append(genesis.RewardMultipliers, types.GenesisRewardMultiplier{
			Validator:        validator,
			Denom:            key.K2(),
			RewardMultiplier: value,
		})
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("ExportGenesis: walk reward multipliers: %w", err)
	}

	err = k.UsersLastRewardMultiplier.Walk(ctx, nil, func(key collections.Triple[sdk.AccAddress, []byte, string], value string) (bool, error) {
		user, err := k.addressCodec.BytesToString(key.K1())
		if err != nil {
			return true, fmt.Errorf("ExportGenesis: user address for last reward multiplier denom %s: %w", key.K3(), err)
		}
		validator, err := k.stakingKeeper.ValidatorAddressCodec().BytesToString(key.K2())
		if err != nil {
			return true, fmt.Errorf("ExportGenesis: validator address for user %s last reward multiplier denom %s: %w", user, key.K3(), err)
		}
		genesis.UserLastRewardMultipliers = append(genesis.UserLastRewardMultipliers, types.GenesisUserLastRewardMultiplier{
			User:             user,
			Validator:        validator,
			Denom:            key.K3(),
			RewardMultiplier: value,
		})
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("ExportGenesis: walk user last reward multipliers: %w", err)
	}

	err = k.LastRewardHandlingTime.Walk(ctx, nil, func(key []byte, value int64) (bool, error) {
		validator, err := k.stakingKeeper.ValidatorAddressCodec().BytesToString(key)
		if err != nil {
			return true, fmt.Errorf("ExportGenesis: validator address for last reward handling time: %w", err)
		}
		genesis.LastRewardHandlingTimes = append(genesis.LastRewardHandlingTimes, types.GenesisLastRewardHandlingTime{
			Validator:              validator,
			LastRewardHandlingTime: value,
		})
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("ExportGenesis: walk last reward handling times: %w", err)
	}
	return nil
}
