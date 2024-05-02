package types

import (
	"fmt"
)

func GrantTokenDenom(expiryHeight uint64) string {
	return fmt.Sprintf("blobgrant/expiry/%d", expiryHeight)
}

func ParseExpiryOfGrantToken(denom string) (uint64, error) {
	var expiryHeight uint64
	_, err := fmt.Sscanf(denom, "blobgrant/expiry/%d", &expiryHeight)
	if err != nil {
		return 0, err
	}
	return expiryHeight, nil
}
