package types

// GetDenom returns the denom for a base YBT token
func GetDenom(creator string) string {
	return ModuleName + "/" + creator
}