package mint

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	minttypes "cosmossdk.io/x/mint/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	liquidityincentivetypes "github.com/sunriselayer/sunrise/x/liquidityincentive/types"

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
func ProvideMintFn(bankKeeper BankKeeper) minttypes.MintFn {
	return func(ctx context.Context, env appmodule.Environment, minter *minttypes.Minter, epochID string, epochNumber int64) error {
		// in this fn we ignore epochNumber as we don't care what epoch we are in, we just assume we are being called every minute.
		if epochID != "minute" {
			return nil
		}

		supplyBond := bankKeeper.GetSupply(ctx, consts.BondDenom)
		supplyFee := bankKeeper.GetSupply(ctx, consts.FeeDenom)
		totalSupply := supplyBond.Amount.Add(supplyFee.Amount)

		annualProvision := CalculateAnnualProvision(
			ctx,
			InflationRateCapInitial,
			InflationRateCapMinimum,
			DisinflationRate,
			SupplyCap,
			Genesis,
			totalSupply,
		)

		// to get a more accurate amount of tokens minted, we get, and later store, last minting time.

		// if this is the first time minting, we initialize the minter.Data with the current time - 60s
		// to mint tokens at the beginning. Note: this is a custom behavior to avoid breaking tests.
		if minter.Data == nil {
			minter.Data = make([]byte, 8)
			binary.BigEndian.PutUint64(minter.Data, (uint64)(env.HeaderService.HeaderInfo(ctx).Time.Unix()-60))
		}

		lastMint := binary.BigEndian.Uint64(minter.Data)
		binary.BigEndian.PutUint64(minter.Data, (uint64)(env.HeaderService.HeaderInfo(ctx).Time.Unix()))

		// calculate the amount of tokens to mint, based on the time since the last mint
		secondsSinceLastMint := env.HeaderService.HeaderInfo(ctx).Time.Unix() - (int64)(lastMint)

		blockProvision := annualProvision.Mul(math.NewInt(secondsSinceLastMint)).Quo(math.NewInt(secondsPerYear))

		if blockProvision.IsPositive() {
			res, err := env.QueryRouterService.Invoke(ctx, &liquidityincentivetypes.QueryParamsRequest{})
			if err != nil {
				return err
			}
			liquidityIncentiveParamsRes, ok := res.(*liquidityincentivetypes.QueryParamsResponse)
			if !ok {
				return fmt.Errorf("unexpected response type: %T", res)
			}
			stakingRewardRatio, err := math.LegacyNewDecFromStr(liquidityIncentiveParamsRes.Params.StakingRewardRatio)
			if err != nil {
				return err
			}
			feeProvision := stakingRewardRatio.MulInt(blockProvision).TruncateInt()
			bondProvision := blockProvision.Sub(feeProvision)

			if feeProvision.IsPositive() {
				feeCoins := sdk.NewCoins(sdk.NewCoin(consts.FeeDenom, feeProvision))

				if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, feeCoins); err != nil {
					return err
				}
				if err := bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, authtypes.FeeCollectorName, feeCoins); err != nil {
					return err
				}
			}

			if bondProvision.IsPositive() {
				bondCoins := sdk.NewCoins(sdk.NewCoin(consts.BondDenom, bondProvision))

				if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, bondCoins); err != nil {
					return err
				}
				if err := bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, authtypes.FeeCollectorName, bondCoins); err != nil {
					return err
				}
			}
		}

		return nil
	}
}
