package mint

import (
	"context"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"

	"github.com/sunriselayer/sunrise/app/consts"
)

var (
	InflationRateCapInitial = math.LegacyMustNewDecFromStr("0.1")
	InflationRateCapMinimum = math.LegacyMustNewDecFromStr("0.02")
	DisinflationRate        = math.LegacyMustNewDecFromStr("0.08")
	SupplyCap               = math.NewInt(1_000_000_000).Mul(math.NewInt(1_000_000))
	Genesis                 = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
)

const (
	secondsPerYear = 31536000
)

type BankKeeper interface {
	GetSupply(ctx context.Context, denom string) sdk.Coin
	MintCoins(ctx context.Context, moduleName string, coins sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule string, recipientModule string, amt sdk.Coins) error
}

// ProvideMintFn returns the function used in x/mint's endblocker to mint new tokens.
// Note that this function can not have the mint keeper as a parameter because it would create a cyclic dependency.
func ProvideMintFn(bankKeeper BankKeeper) mintkeeper.MintFn {
	return func(ctx sdk.Context, k *mintkeeper.Keeper) error {
		// fetch stored minter & params
		minter, err := k.Minter.Get(ctx)
		if err != nil {
			return err
		}

		params, err := k.Params.Get(ctx)
		if err != nil {
			return err
		}

		supplyBond := bankKeeper.GetSupply(ctx, consts.BondDenom)
		supplyMint := bankKeeper.GetSupply(ctx, consts.MintDenom)
		totalSupply := supplyBond.Amount.Add(supplyMint.Amount)

		annualProvisions := CalculateAnnualProvision(
			ctx,
			InflationRateCapInitial,
			InflationRateCapMinimum,
			DisinflationRate,
			SupplyCap,
			Genesis,
			totalSupply,
		)
		inflationRate := annualProvisions.QuoInt(totalSupply)

		minter.Inflation = inflationRate
		minter.AnnualProvisions = annualProvisions
		if err = k.Minter.Set(ctx, minter); err != nil {
			return err
		}

		blockProvisions := annualProvisions.QuoInt(math.NewInt(int64(params.BlocksPerYear))).TruncateInt()

		if blockProvisions.IsPositive() {
			mintCoins := sdk.NewCoins(sdk.NewCoin(consts.MintDenom, blockProvisions))

			if err := k.MintCoins(ctx, mintCoins); err != nil {
				return err
			}
			if err := k.AddCollectedFees(ctx, mintCoins); err != nil {
				return err
			}
		}
		return nil
	}
}
