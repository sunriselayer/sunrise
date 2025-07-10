package keeper_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
	"go.uber.org/mock/gomock"
)

func TestKeeper_ValidateLastRewardHandlingTime(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	validatorAddr := sdk.ValAddress([]byte("validator"))
	params := types.DefaultParams()
	params.RewardPeriod = time.Hour
	require.NoError(f.keeper.Params.Set(f.ctx, params))

	// Case 1: First time handling
	err := f.keeper.ValidateLastRewardHandlingTime(f.ctx, validatorAddr)
	require.NoError(err)

	// Case 2: Within reward period
	err = f.keeper.ValidateLastRewardHandlingTime(f.ctx, validatorAddr)
	require.NoError(err)

	// Case 3: After reward period
	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	sdkCtx = sdkCtx.WithBlockTime(sdkCtx.BlockTime().Add(params.RewardPeriod))
	f.ctx = sdk.WrapSDKContext(sdkCtx)

	err = f.keeper.ValidateLastRewardHandlingTime(f.ctx, validatorAddr)
	require.NoError(err)
}

func TestKeeper_HandleModuleAccountRewards(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	valAddressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())

	moduleAddr := sdk.AccAddress([]byte("module"))
	valAddr1 := sdk.ValAddress([]byte("validator1"))
	valAddr2 := sdk.ValAddress([]byte("validator2"))

	delegations := []stakingtypes.Delegation{
		{ValidatorAddress: valAddr1.String()},
		{ValidatorAddress: valAddr2.String()},
	}

	f.mocks.AccountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr).AnyTimes()
	f.mocks.StakingKeeper.EXPECT().IterateDelegatorDelegations(gomock.Any(), moduleAddr, gomock.Any()).
		DoAndReturn(func(ctx context.Context, delegator sdk.AccAddress, cb func(delegation stakingtypes.Delegation) (stop bool)) error {
			for _, del := range delegations {
				if cb(del) {
					break
				}
			}
			return nil
		})

	// Mock HandleModuleAccountRewardsByValidator for each validator
	f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec).AnyTimes()
	f.mocks.DistributionKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), moduleAddr, valAddr1).Return(sdk.NewCoins(), nil)
	f.mocks.DistributionKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), moduleAddr, valAddr2).Return(sdk.NewCoins(), nil)

	err := f.keeper.HandleModuleAccountRewards(f.ctx)
	require.NoError(err)
}

func TestKeeper_HandleModuleAccountRewardsByValidator(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)
	valAddressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())

	valAddr := sdk.ValAddress([]byte("validator"))
	valAddrStr := valAddr.String()

	moduleAddr := sdk.AccAddress([]byte("module"))
	rewardSaverAddr := types.RewardSaverAddress(valAddrStr)
	rewards := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))
	totalShare := math.NewInt(1000)

	// Setup mocks
	f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
	f.mocks.AccountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr)
	f.mocks.DistributionKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), moduleAddr, valAddr).Return(rewards, nil)
	f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, rewardSaverAddr, rewards).Return(nil)
	f.mocks.BankKeeper.EXPECT().GetSupply(gomock.Any(), types.NonVotingShareTokenDenom(valAddrStr)).Return(sdk.NewCoin(types.NonVotingShareTokenDenom(valAddrStr), totalShare))

	// Set params for ValidateLastRewardHandlingTime
	params := types.DefaultParams()
	params.RewardPeriod = time.Hour
	require.NoError(f.keeper.Params.Set(f.ctx, params))

	err := f.keeper.HandleModuleAccountRewardsByValidator(f.ctx, valAddrStr)
	require.NoError(err)
}
