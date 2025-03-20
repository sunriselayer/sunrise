package types

import (
	"fmt"

	"cosmossdk.io/math"
)

func LiquidStakingTokenDenom(validatorAddress string) string {
	return fmt.Sprintf("%s/%s", ModuleName, validatorAddress)
}

func RewardSaverModuleAccount(validatorAddress string) string {
	return fmt.Sprintf("%s/reward_saver/%s", ModuleName, validatorAddress)
}

func CalculateLiquidUnstakeOutputAmount(stakedAmount, lstSupplyOld, lstAmount math.Int) (math.Int, error) {
	stakedAmountDec, err := math.NewDecFromString(stakedAmount.String())
	if err != nil {
		return math.Int{}, err
	}

	lstSupplyOldDec, err := math.NewDecFromString(lstSupplyOld.String())
	if err != nil {
		return math.Int{}, err
	}

	ratio, err := stakedAmountDec.Quo(lstSupplyOldDec)
	if err != nil {
		return math.Int{}, err
	}

	lstAmountDec, err := math.NewDecFromString(lstAmount.String())
	if err != nil {
		return math.Int{}, err
	}

	outputAmountDec, err := lstAmountDec.Mul(ratio)
	if err != nil {
		return math.Int{}, err
	}
	outputAmount, err := outputAmountDec.SdkIntTrim()
	if err != nil {
		return math.Int{}, err
	}

	return outputAmount, nil
}

func CalculateRewardMultiplierNew(rewardMultiplierOld math.Dec, rewardAmount math.Int, lstSupplyNew math.Int) (math.Dec, error) {
	multiplierDiffNumerator, err := math.NewDecFromString(rewardAmount.String())
	if err != nil {
		return math.Dec{}, err
	}
	multiplierDiffDenominator, err := math.NewDecFromString(lstSupplyNew.String())
	if err != nil {
		return math.Dec{}, err
	}
	multiplierDiff, err := multiplierDiffNumerator.Quo(multiplierDiffDenominator)
	if err != nil {
		return math.Dec{}, err
	}

	multiplierNew, err := rewardMultiplierOld.Add(multiplierDiff)
	if err != nil {
		return math.Dec{}, err
	}

	return multiplierNew, nil
}

func CalculateSlashingCompensation(stakedAmount, lstSupplyOld, feeCoinRewardAmount math.Int) (compensation math.Int, distribution math.Int) {
	slashed := lstSupplyOld.Sub(stakedAmount)

	if feeCoinRewardAmount.GT(slashed) {
		compensation = slashed
		distribution = feeCoinRewardAmount.Sub(compensation)
	} else {
		compensation = feeCoinRewardAmount
		distribution = math.NewInt(0)
	}

	return
}
