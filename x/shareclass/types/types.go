package types

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"

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
	if totalBonded.IsZero() {
		return amount, nil
	}

	amountDec, err := math.NewDecFromString(amount.String())
	if err != nil {
		return math.Int{}, err
	}

	totalShareDec, err := math.NewDecFromString(totalShare.String())
	if err != nil {
		return math.Int{}, err
	}

	totalBondedDec, err := math.NewDecFromString(totalBonded.String())
	if err != nil {
		return math.Int{}, err
	}

	ratio, err := amountDec.Quo(totalBondedDec)
	if err != nil {
		return math.Int{}, err
	}

	shareDec, err := ratio.Mul(totalShareDec)
	if err != nil {
		return math.Int{}, err
	}

	return shareDec.SdkIntTrim()
}

func CalculateAmountByShare(totalShare, totalBonded, share math.Int) (math.Int, error) {
	if totalShare.IsZero() {
		return share, nil
	}

	shareDec, err := math.NewDecFromString(share.String())
	if err != nil {
		return math.Int{}, err
	}

	totalShareDec, err := math.NewDecFromString(totalShare.String())
	if err != nil {
		return math.Int{}, err
	}

	totalBondedDec, err := math.NewDecFromString(totalBonded.String())
	if err != nil {
		return math.Int{}, err
	}

	ratio, err := shareDec.Quo(totalShareDec)
	if err != nil {
		return math.Int{}, err
	}

	outputAmountDec, err := totalBondedDec.Mul(ratio)
	if err != nil {
		return math.Int{}, err
	}

	return outputAmountDec.SdkIntTrim()
}

// CalculateReward calculates the reward for a user
// reward = (rewardMultiplier - userLastRewardMultiplier) * principal
func CalculateReward(rewardMultiplier, userLastRewardMultiplier math.Dec, share math.Int) (math.Int, error) {
	shareDec, err := math.NewDecFromString(share.String())
	if err != nil {
		return math.Int{}, err
	}

	multiplierDiff, err := rewardMultiplier.Sub(userLastRewardMultiplier)
	if err != nil {
		return math.Int{}, err
	}

	rewardAmountDec, err := multiplierDiff.Mul(shareDec)
	if err != nil {
		return math.Int{}, err
	}
	rewardAmount, err := rewardAmountDec.SdkIntTrim()
	if err != nil {
		return math.Int{}, err
	}

	return rewardAmount, nil
}

func CalculateRewardMultiplierNew(rewardMultiplierOld math.Dec, rewardAmount math.Int, totalShare math.Int) (math.Dec, error) {
	multiplierDiffNumerator, err := math.NewDecFromString(rewardAmount.String())
	if err != nil {
		return math.Dec{}, err
	}
	multiplierDiffDenominator, err := math.NewDecFromString(totalShare.String())
	if err != nil {
		return math.Dec{}, err
	}
	multiplierDiff, err := multiplierDiffNumerator.Quo(multiplierDiffDenominator)
	if err != nil {
		return math.Dec{}, err
	}

	multiplierNew, err := rewardMultiplierOld.Add(multiplierDiff)
	if err != nil {
		return math.Dec{}, err
	}

	return multiplierNew, nil
}
