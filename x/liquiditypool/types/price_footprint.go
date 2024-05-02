package types

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"time"
)

type AscendingOrderedPriceFootprints []PriceFootprint

func CalculateTwap(footprints AscendingOrderedPriceFootprints) (*math.LegacyDec, error) {
	if len(footprints) <= 1 {
		return nil, errors.Wrapf(ErrInsufficientFootprintForTwap, "length of price footprints %d is shorter than 2", len(footprints))
	}

	var sum math.LegacyDec

	for i := 1; i < len(footprints); i++ {
		footprint := footprints[i]
		previousFootprint := footprints[i-1]

		duration := footprint.Timestamp.Sub(previousFootprint.Timestamp)
		durationNanoSecInt := math.NewInt(duration.Nanoseconds())
		weightedPrice := footprint.Price.Mul(durationNanoSecInt.ToLegacyDec())

		sum = sum.Add(weightedPrice)
	}

	totalDuration := footprints[len(footprints)-1].Timestamp.Sub(footprints[0].Timestamp)
	totalDurationNanoSecInt := math.NewInt(totalDuration.Nanoseconds())
	twap := sum.Quo(totalDurationNanoSecInt.ToLegacyDec())

	return &twap, nil
}

func NotDeprecatedTimestampStartIndex(footprints AscendingOrderedPriceFootprints, now time.Time, period time.Duration) uint64 {
	deprecatedTimestamp := now.Add(-period)

	for i, footprint := range footprints {
		if footprint.Timestamp.After(deprecatedTimestamp) {
			return uint64(i)
		}
	}

	return uint64(len(footprints))
}
