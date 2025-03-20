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
