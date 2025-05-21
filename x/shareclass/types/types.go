package types

import (
	"fmt"
	"regexp"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func NonVotingShareTokenDenom(validatorAddress string) string {
	return fmt.Sprintf("%s/non-voting-share/%s", ModuleName, validatorAddress)
}

func NonVotingShareTokenDenomRegexp() *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s/non-voting-share/([^/]+)$", ModuleName))
}

func RewardSaverAddress(validatorAddress string) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("%s/reward_saver/%s", ModuleName, validatorAddress))
}

func CalculateShareByAmount(totalShare, totalBonded, amount math.Int) (math.Int, error) {
	if totalShare.IsZero() {
		return amount, nil
	}
	if totalBonded.IsZero() {
		return amount, nil
	}

	amountDec, err := math.LegacyNewDecFromStr(amount.String())
	if err != nil {
		return math.Int{}, err
	}

	totalShareDec, err := math.LegacyNewDecFromStr(totalShare.String())
	if err != nil {
		return math.Int{}, err
	}

	totalBondedDec, err := math.LegacyNewDecFromStr(totalBonded.String())
	if err != nil {
		return math.Int{}, err
	}

	// totalBonded is not zero
	ratio := amountDec.Quo(totalBondedDec)

	shareDec := ratio.Mul(totalShareDec)

	return shareDec.TruncateInt(), nil
}

func CalculateAmountByShare(totalShare, totalBonded, share math.Int) (math.Int, error) {
	if totalShare.IsZero() {
		return share, nil
	}

	shareDec, err := math.LegacyNewDecFromStr(share.String())
	if err != nil {
		return math.Int{}, err
	}

	totalShareDec, err := math.LegacyNewDecFromStr(totalShare.String())
	if err != nil {
		return math.Int{}, err
	}

	totalBondedDec, err := math.LegacyNewDecFromStr(totalBonded.String())
	if err != nil {
		return math.Int{}, err
	}

	// totalShare is not zero
	ratio := shareDec.Quo(totalShareDec)

	outputAmountDec := totalBondedDec.Mul(ratio)

	return outputAmountDec.TruncateInt(), nil
}

// CalculateReward calculates the reward for a user
// reward = (rewardMultiplier - userLastRewardMultiplier) * principal
func CalculateReward(rewardMultiplier, userLastRewardMultiplier math.LegacyDec, share math.Int) (math.Int, error) {
	shareDec, err := math.LegacyNewDecFromStr(share.String())
	if err != nil {
		return math.Int{}, err
	}

	multiplierDiff := rewardMultiplier.Sub(userLastRewardMultiplier)
	rewardAmountDec := multiplierDiff.Mul(shareDec)
	rewardAmount := rewardAmountDec.TruncateInt()

	return rewardAmount, nil
}

func CalculateRewardMultiplierNew(rewardMultiplierOld math.LegacyDec, rewardAmount math.Int, totalShare math.Int) (math.LegacyDec, error) {
	multiplierDiffNumerator, err := math.LegacyNewDecFromStr(rewardAmount.String())
	if err != nil {
		return math.LegacyDec{}, err
	}
	multiplierDiffDenominator, err := math.LegacyNewDecFromStr(totalShare.String())
	if err != nil {
		return math.LegacyDec{}, err
	}
	if multiplierDiffDenominator.IsZero() {
		return math.LegacyDec{}, fmt.Errorf("total share is zero")
	}
	multiplierDiff := multiplierDiffNumerator.Quo(multiplierDiffDenominator)

	multiplierNew := rewardMultiplierOld.Add(multiplierDiff)

	return multiplierNew, nil
}
