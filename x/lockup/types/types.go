package types

import (
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
)

func LockupAccountModule(owner string) string {
	return fmt.Sprintf("%s/%s", ModuleName, owner)
}

func CalculateUnlockedAmount(lockupAmount math.Int, startTime time.Time, endTime time.Time, now time.Time) (math.Int, error) {
	if startTime.After(endTime) {
		return math.Int{}, errorsmod.Wrap(ErrInvalidTimeRange, "start time is after end time")
	}

	if now.Before(startTime) {
		return math.NewInt(0), nil
	}

	if now.After(endTime) {
		return lockupAmount, nil
	}

	numerator := now.Sub(startTime).Milliseconds()
	denominator := endTime.Sub(startTime).Milliseconds()

	numeratorDec := math.NewDecFromInt64(numerator)
	denominatorDec := math.NewDecFromInt64(denominator)

	ratio, err := numeratorDec.Quo(denominatorDec)
	if err != nil {
		return math.Int{}, err
	}

	lockupAmountDec, err := math.NewDecFromString(lockupAmount.String())
	if err != nil {
		return math.Int{}, err
	}

	unlockedAmountDec, err := lockupAmountDec.Mul(ratio)
	if err != nil {
		return math.Int{}, err
	}

	unlockedAmount, err := unlockedAmountDec.SdkIntTrim()
	if err != nil {
		return math.Int{}, err
	}

	return unlockedAmount, nil
}

func CalculateRequiredBalance(lockupAmount, unlockedAmount math.Int) math.Int {
	return lockupAmount.Sub(unlockedAmount)
}

func SendCondition(lockupAmount, unlockedAmount, balance, sendAmount math.Int) bool {
	return balance.Sub(sendAmount).GTE(CalculateRequiredBalance(lockupAmount, unlockedAmount))
}
