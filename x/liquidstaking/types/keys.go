package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "liquidstaking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_liquidstaking"

	// ModuleAccountName is the module account's name
	ModuleAccountName = ModuleName

	DefaultDerivativeDenom = "bsr"

	DenomSeparator = "-"
)

var (
	ParamsKey = []byte("p_liquidstaking")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetLiquidStakingTokenDenom(bondDenom string, valAddr sdk.ValAddress) string {
	return fmt.Sprintf("%s%s%s", bondDenom, DenomSeparator, valAddr.String())
}

// ParseLiquidStakingTokenDenom extracts a validator address from a derivative denom.
func ParseLiquidStakingTokenDenom(denom string) (sdk.ValAddress, error) {
	elements := strings.Split(denom, DenomSeparator)
	if len(elements) != 2 {
		return nil, fmt.Errorf("cannot parse denom %s", denom)
	}

	if elements[0] != DefaultDerivativeDenom {
		return nil, fmt.Errorf("invalid denom prefix, expected %s, got %s", DefaultDerivativeDenom, elements[0])
	}

	addr, err := sdk.ValAddressFromBech32(elements[1])
	if err != nil {
		return nil, fmt.Errorf("invalid denom validator address: %w", err)
	}

	return addr, nil
}
