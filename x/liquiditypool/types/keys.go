package types

import "cosmossdk.io/collections"

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "liquiditypool"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_liquiditypool")

// TODO: use collections

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	PoolKey      = "Pool/value/"
	PoolCountKey = "Pool/count/"
)

const (
	PositionKey       = "Position/value/"
	PositionCountKey  = "Position/count/"
	PositionByPool    = "PositionByPool/"
	PositionByAddress = "PositionByAddress/"
)

const (
	TickInfoKey                  = "TickInfo/value/"
	TickNegativePrefix           = "N"
	TickPositivePrefix           = "P"
	FeePositionAccumulatorPrefix = "FeePositionAccumulator/value/"
	KeyFeePoolAccumulatorPrefix  = "FeePoolAccumulator/value/"
	KeyAccumPrefix               = "Accumulator/Acc/value/"
	KeyAccumulatorPositionPrefix = "Accumulator/Pos/value/"
	KeySeparator                 = "||"
)

func TickIndexFromBytes(bz []byte) (int64, error) {
	if len(bz) != 9 {
		return 0, ErrInvalidTickIndexEncoding
	}

	i := int64(sdk.BigEndianToUint64(bz[1:]))
	if bz[0] == TickNegativePrefix[0] && i >= 0 {
		return 0, ErrInvalidTickIndexEncoding
	} else if bz[0] == TickPositivePrefix[0] && i < 0 {
		return 0, ErrInvalidTickIndexEncoding
	}
	return i, nil
}

func TickIndexToBytes(tickIndex int64) []byte {
	key := make([]byte, 9)
	if tickIndex < 0 {
		copy(key[:1], TickNegativePrefix)
		copy(key[1:], sdk.Uint64ToBigEndian(uint64(tickIndex)))
	} else {
		copy(key[:1], TickPositivePrefix)
		copy(key[1:], sdk.Uint64ToBigEndian(uint64(tickIndex)))
	}

	return key
}

func GetTickInfoIDBytes(poolId uint64, tickIndex int64) []byte {
	bz := KeyTickPrefixByPoolId(poolId)
	bz = append(bz, TickIndexToBytes(tickIndex)...)
	return bz
}

func KeyTickPrefixByPoolId(poolId uint64) []byte {
	bz := KeyPrefix(TickInfoKey)
	bz = append(bz, []byte("/")...)
	bz = binary.BigEndian.AppendUint64(bz, poolId)
	return bz
}

func KeyFeePositionAccumulator(positionId uint64) string {
	return strings.Join([]string{string(FeePositionAccumulatorPrefix), strconv.FormatUint(positionId, 10)}, "|")
}

// This is guaranteed to not contain "||" so it can be used as an accumulator name.
func KeyFeePoolAccumulator(poolId uint64) string {
	poolIdStr := strconv.FormatUint(poolId, 10)
	return strings.Join([]string{string(KeyFeePoolAccumulatorPrefix), poolIdStr}, "/")
}

func FormatKeyAccumPrefix(accumName string) []byte {
	return []byte(fmt.Sprintf(KeyAccumPrefix+"%s", accumName))
}

func FormatKeyAccumulatorPositionPrefix(accumName, name string) []byte {
	return []byte(fmt.Sprintf(KeyAccumulatorPositionPrefix+"%s"+KeySeparator+"%s", accumName, name))
}

func PositionByPoolPrefix(poolId uint64) []byte {
	return append([]byte(PositionByPool), sdk.Uint64ToBigEndian(poolId)...)
}

func PositionByAddressPrefix(addr string) []byte {
	return append([]byte(PositionByAddress), addr...)
}
