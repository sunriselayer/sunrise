package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetBaseYbtTokenDenom returns the denom for a base YBT token
func GetBaseYbtTokenDenom(baseYbtCreator string) string {
	return "ybtbase/" + baseYbtCreator
}

// GetBaseYbtYieldPoolAddress returns the module account address for a base YBT token's yield pool
func GetBaseYbtYieldPoolAddress(baseYbtCreator string) sdk.AccAddress {
	return authtypes.NewModuleAddress("ybtbase/yield/" + baseYbtCreator)
}