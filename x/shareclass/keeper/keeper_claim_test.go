package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
	"go.uber.org/mock/gomock"
)

func TestKeeper_RewardMultiplier(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	validatorAddr := sdk.ValAddress([]byte("validator"))
	denom := "test"
	multiplier := math.LegacyNewDec(2)

	// Set
	err := f.keeper.SetRewardMultiplier(f.ctx, validatorAddr, denom, multiplier)
	require.NoError(err)

	// Get
	result, err := f.keeper.GetRewardMultiplier(f.ctx, validatorAddr, denom)
	require.NoError(err)
	require.Equal(multiplier, result)

	// Get non-existent
	result, err = f.keeper.GetRewardMultiplier(f.ctx, sdk.ValAddress([]byte("non-existent")), denom)
	require.NoError(err)
	require.True(result.IsZero())
}

func TestKeeper_UsersLastRewardMultiplier(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	sender := sdk.AccAddress([]byte("sender"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	denom := "test"
	multiplier := math.LegacyNewDec(2)

	// Set
	err := f.keeper.SetUserLastRewardMultiplier(f.ctx, sender, validatorAddr, denom, multiplier)
	require.NoError(err)

	// Get
	result, err := f.keeper.GetUserLastRewardMultiplier(f.ctx, sender, validatorAddr, denom)
	require.NoError(err)
	require.Equal(multiplier, result)

	// Get non-existent
	result, err = f.keeper.GetUserLastRewardMultiplier(f.ctx, sdk.AccAddress([]byte("non-existent")), validatorAddr, denom)
	require.NoError(err)
	require.True(result.IsZero())
}

func TestKeeper_GetClaimableRewards(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	sender := sdk.AccAddress([]byte("sender"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	denom := "test"
	rewardSaverAddr := types.RewardSaverAddress(validatorAddr.String())

	// Setup mocks
	f.mocks.BankKeeper.EXPECT().GetAllBalances(gomock.Any(), rewardSaverAddr).Return(sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(100))))
	f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin(denom, math.NewInt(100)))

	// Set multipliers
	require.NoError(f.keeper.SetRewardMultiplier(f.ctx, validatorAddr, denom, math.LegacyNewDec(2)))
	require.NoError(f.keeper.SetUserLastRewardMultiplier(f.ctx, sender, validatorAddr, denom, math.LegacyNewDec(1)))

	rewards, err := f.keeper.GetClaimableRewards(f.ctx, sender, validatorAddr)
	require.NoError(err)
	require.Equal(sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(100))), rewards)
}

func TestKeeper_ClaimRewards(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	sender := sdk.AccAddress([]byte("sender"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	denom := "test"
	rewardSaverAddr := types.RewardSaverAddress(validatorAddr.String())
	expectedRewards := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(100)))

	// Setup mocks
	f.mocks.BankKeeper.EXPECT().GetAllBalances(gomock.Any(), rewardSaverAddr).Return(expectedRewards)
	f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin(denom, math.NewInt(100)))
	f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), rewardSaverAddr, sender, expectedRewards).Return(nil)

	// Set multipliers
	require.NoError(f.keeper.SetRewardMultiplier(f.ctx, validatorAddr, denom, math.LegacyNewDec(2)))
	require.NoError(f.keeper.SetUserLastRewardMultiplier(f.ctx, sender, validatorAddr, denom, math.LegacyNewDec(1)))

	rewards, err := f.keeper.ClaimRewards(f.ctx, sender, validatorAddr)
	require.NoError(err)
	require.Equal(expectedRewards, rewards)
}

func TestKeeper_GetClaimableRewardsByDenom_NoStateChange(t *testing.T) {
	f := initFixture(t)
	require := require.New(t)

	sender := sdk.AccAddress([]byte("sender"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	denom := "test"
	initialUserMultiplier := math.LegacyNewDec(1)
	globalMultiplier := math.LegacyNewDec(2)

	// Setup initial state
	require.NoError(f.keeper.SetRewardMultiplier(f.ctx, validatorAddr, denom, globalMultiplier))
	require.NoError(f.keeper.SetUserLastRewardMultiplier(f.ctx, sender, validatorAddr, denom, initialUserMultiplier))

	// Mock GetShare to return a value
	shareAmount := math.NewInt(100)
	f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), sender, types.NonVotingShareTokenDenom(validatorAddr.String())).Return(sdk.NewCoin(denom, shareAmount))

	// This call should NOT update the user's last reward multiplier
	_, err := f.keeper.GetClaimableRewardsByDenom(f.ctx, sender, validatorAddr, denom)
	require.NoError(err)

	// Verify that the user's last reward multiplier has NOT been updated
	updatedUserMultiplier, err := f.keeper.GetUserLastRewardMultiplier(f.ctx, sender, validatorAddr, denom)
	require.NoError(err)
	require.Equal(initialUserMultiplier, updatedUserMultiplier, "UserLastRewardMultiplier should not be updated on a Get call")
}
