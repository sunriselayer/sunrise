package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunrise-zone/sunrise-app/testutil"
	"github.com/sunrise-zone/sunrise-app/x/liquidstaking/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func (suite *KeeperTestSuite) TestCollectStakingRewards() {
	_, addrs := testutil.GeneratePrivKeyAddressPairs(5)
	valAccAddr1, delegator := addrs[0], addrs[1]
	valAddr1 := sdk.ValAddress(valAccAddr1)

	initialBalance := i(1e9)
	delegateAmount := i(100e6)

	suite.NoError(suite.FundModuleAccount(
		suite.Ctx,
		distrtypes.ModuleName,
		sdk.NewCoins(
			sdk.NewCoin("usr", initialBalance),
		),
	))

	suite.CreateAccountWithAddress(valAccAddr1, suite.NewBondCoins(initialBalance))
	suite.CreateAccountWithAddress(delegator, suite.NewBondCoins(initialBalance))

	suite.CreateNewUnbondedValidator(valAddr1, initialBalance)
	suite.CreateDelegation(valAddr1, delegator, delegateAmount)
	suite.StakingKeeper.EndBlocker(suite.Ctx)

	// Transfers delegation to module account
	_, err := suite.Keeper.MintDerivative(suite.Ctx, delegator, valAddr1, suite.NewBondCoin(delegateAmount))
	suite.Require().NoError(err)

	validator, err := suite.StakingKeeper.GetValidator(suite.Ctx, valAddr1)
	suite.Require().NoError(err)

	suite.Ctx = suite.Ctx.WithBlockHeight(2)

	distrKeeper := suite.App.DistrKeeper
	stakingKeeper := suite.App.StakingKeeper
	accKeeper := suite.App.AccountKeeper
	liquidMacc := accKeeper.GetModuleAccount(suite.Ctx, types.ModuleAccountName)

	// Add rewards
	rewardCoins := sdk.NewDecCoins(sdk.NewDecCoin("usr", sdkmath.NewInt(500e6)))
	distrKeeper.AllocateTokensToValidator(suite.Ctx, validator, rewardCoins)

	delegation, err := stakingKeeper.GetDelegation(suite.Ctx, liquidMacc.GetAddress(), valAddr1)
	suite.Require().NoError(err)

	// Get amount of rewards
	endingPeriod, _ := distrKeeper.IncrementValidatorPeriod(suite.Ctx, validator)
	delegationRewards, _ := distrKeeper.CalculateDelegationRewards(suite.Ctx, validator, delegation, endingPeriod)
	truncatedRewards, _ := delegationRewards.TruncateDecimal()

	suite.Run("collect staking rewards", func() {
		// Collect rewards
		derivativeDenom := suite.Keeper.GetLiquidStakingTokenDenom(valAddr1)
		rewards, err := suite.Keeper.CollectStakingRewardsByDenom(suite.Ctx, derivativeDenom, types.ModuleName)
		suite.Require().NoError(err)
		suite.Require().Equal(truncatedRewards, rewards)

		suite.True(rewards.AmountOf("usr").IsPositive())

		// Check balances
		suite.AccountBalanceEqual(liquidMacc.GetAddress(), rewards)
	})

	suite.Run("collect staking rewards with non-validator", func() {
		// acc2 not a validator
		derivativeDenom := suite.Keeper.GetLiquidStakingTokenDenom(sdk.ValAddress(addrs[2]))
		_, err := suite.Keeper.CollectStakingRewardsByDenom(suite.Ctx, derivativeDenom, types.ModuleName)
		suite.Require().Error(err)
		suite.Require().Equal("validator does not exist", err.Error())
	})

	suite.Run("collect staking rewards with invalid denom", func() {
		derivativeDenom := "bstake"
		_, err := suite.Keeper.CollectStakingRewardsByDenom(suite.Ctx, derivativeDenom, types.ModuleName)
		suite.Require().Error(err)
		suite.Require().Equal("cannot parse denom bstake", err.Error())
	})
}
