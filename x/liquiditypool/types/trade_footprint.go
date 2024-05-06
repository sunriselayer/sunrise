package types

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
)

func CalculateTwap(footprints []TradeFootprint) (*math.LegacyDec, error) {
	if len(footprints) == 0 {
		return nil, errors.Wrap(ErrInsufficientFootprintForTwap, "empty footprint")
	}

	var sumPrice math.LegacyDec
	var sumVolume math.Int

	for _, footprint := range footprints {
		sumPrice = sumPrice.Add(footprint.Price)
		sumVolume = sumVolume.Add(footprint.Volume)
	}

	twap := sumPrice.Quo(math.LegacyNewDecFromInt(sumVolume))

	return &twap, nil
}
