package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
)

func LockupAccountModule(owner string) string {
	return fmt.Sprintf("%s/%s", ModuleName, owner)
}

func (lockup LockupAccount) GetLockCoinInfo(blockTime int64) (unlockedAmount, lockedAmount math.Int, err error) {
	startTime := lockup.StartTime
	endTime := lockup.EndTime
	originalLocking := lockup.OriginalLocking
	additionalLocking := lockup.AdditionalLocking
	totalLocking := originalLocking.Add(additionalLocking)

	if startTime > endTime {
		return math.Int{}, math.Int{}, errorsmod.Wrap(ErrInvalidTimeRange, "start time is after end time")
	}

	if blockTime < startTime {
		return math.NewInt(0), totalLocking, nil
	}

	if blockTime > endTime {
		return totalLocking, math.NewInt(0), nil
	}

	// calculate the locking scalar
	x := blockTime - startTime
	y := endTime - startTime
	s := math.LegacyNewDec(x).Quo(math.LegacyNewDec(y))

	unlockedAmt := math.LegacyNewDecFromInt(totalLocking).Mul(s).RoundInt()
	lockedAmt := totalLocking.Sub(unlockedAmt)

	return unlockedAmt, lockedAmt, nil
}

func (lockup LockupAccount) GetNotBondedLockedAmount(lockedAmount math.Int) math.Int {
	delegatedLockingAmt := lockup.DelegatedLocking
	x := math.MinInt(lockedAmount, delegatedLockingAmt)
	lockedAmt := lockedAmount.Sub(x)
	return lockedAmt
}
