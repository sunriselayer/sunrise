package types

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	"cosmossdk.io/core/address"
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

var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("params/")

	PoolsKeyPrefix              = collections.NewPrefix("pools/")
	PoolIdKey                   = collections.NewPrefix("pool_id/")
	PositionsKeyPrefix          = collections.NewPrefix("positions/")
	PositionIdKey               = collections.NewPrefix("position_id/")
	PositionsPoolIdIndexPrefix  = collections.NewPrefix("positions_by_pool_id/")
	PositionsAddressIndexPrefix = collections.NewPrefix("positions_by_address/")
)

type PositionsIndexes struct {
	PoolId  *indexes.Multi[uint64, uint64, Position]
	Address *indexes.Multi[sdk.AccAddress, uint64, Position]
}

func (i PositionsIndexes) IndexesList() []collections.Index[uint64, Position] {
	return []collections.Index[uint64, Position]{
		i.PoolId,
		i.Address,
	}
}

func NewPositionsIndexes(sb *collections.SchemaBuilder, addressCodec address.Codec) PositionsIndexes {
	return PositionsIndexes{
		PoolId: indexes.NewMulti(
			sb,
			PositionsPoolIdIndexPrefix,
			"positions_by_pool_id",
			collections.Uint64Key,
			collections.Uint64Key,
			func(_ uint64, v Position) (uint64, error) {
				return v.PoolId, nil
			},
		),
		Address: indexes.NewMulti(
			sb,
			PositionsAddressIndexPrefix,
			"positions_by_address",
			sdk.AccAddressKey,
			collections.Uint64Key,
			func(_ uint64, v Position) (sdk.AccAddress, error) {
				return addressCodec.StringToBytes(v.Address)
			},
		),
	}
}

var (
	PoolsKeyCodec     = collections.Uint64Key
	PositionsKeyCodec = collections.Uint64Key
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	TickInfoKey                  = "tick_info/"
	TickNegativePrefix           = "N"
	TickPositivePrefix           = "P"
	FeePositionAccumulatorPrefix = "fee_position_accumulator/"
	KeyFeePoolAccumulatorPrefix  = "fee_pool_accumulator/"
	KeyAccumPrefix               = "accumulator/"
	KeyAccumulatorPositionPrefix = "accumulator_position/"
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
	bz = binary.BigEndian.AppendUint64(bz, poolId)
	return bz
}

func KeyFeePositionAccumulator(positionId uint64) string {
	return strings.Join([]string{string(FeePositionAccumulatorPrefix), strconv.FormatUint(positionId, 10)}, "|")
}

// This is guaranteed to not contain "||" so it can be used as an accumulator name.
func KeyFeePoolAccumulator(poolId uint64) string {
	poolIdStr := strconv.FormatUint(poolId, 10)
	return KeyFeePoolAccumulatorPrefix + poolIdStr
}

func FormatKeyAccumPrefix(accumName string) []byte {
	return fmt.Appendf(nil, KeyAccumPrefix+"%s", accumName)
}

func FormatKeyAccumulatorPositionPrefix(accumName, name string) []byte {
	return fmt.Appendf(nil, KeyAccumulatorPositionPrefix+"%s"+KeySeparator+"%s", accumName, name)
}
