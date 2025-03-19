package types

import (
	"fmt"
)

func LiquidStakingTokenDenom(validatorAddress string) string {
	return fmt.Sprintf("%s/%s", ModuleName, validatorAddress)
}
