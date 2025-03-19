package types

import (
	"fmt"
)

func LiquidStakingTokenDenom(validatorAddress string) string {
	return fmt.Sprintf("%s/%s", ModuleName, validatorAddress)
}

func RewardSaverModuleAccount() string {
	return fmt.Sprintf("%s/reward_saver", ModuleName)
}
