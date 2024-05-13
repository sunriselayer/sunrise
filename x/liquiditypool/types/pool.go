package types

import (
	"fmt"

	"cosmossdk.io/math"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquiditypool/cfmm"
)

func LpTokenDenom(poolId uint64) string {
	return fmt.Sprintf("%s/%d", ModuleName, poolId)
}

func LpTokenValueInQuoteUnit(baseAmount math.Int, quoteAmount math.Int, price math.LegacyDec) math.LegacyDec {
	baseAmountQuotedValue := price.MulInt(baseAmount)
	quoteAmountQuotedValue := quoteAmount.ToLegacyDec()

	return baseAmountQuotedValue.Add(quoteAmountQuotedValue)
}

func PoolModuleName(poolId uint64) string {
	return fmt.Sprintf("%s/%d", ModuleName, poolId)
}

func PoolTreasuryModuleName(poolId uint64) string {
	return fmt.Sprintf("%s/%d/treasury", ModuleName, poolId)
}

func UnpackCfmm(ctx sdk.Context, pool Pool) (cfmm.ConstantFunctionMarketMaker, error) {
	var cfmm cfmm.ConstantFunctionMarketMaker
	err := ctx.UnpackAny(&pool.Cfmm, cfmm)

	return cfmm, err
}

func CalculateX(y math.Int, k math.LegacyDec, cfmm cfmm.ConstantFunctionMarketMaker) (*math.Int, error) {
	return nil, nil
}

func CalculateDx(y math.Int, dy math.Int, k math.LegacyDec, cfmm cfmm.ConstantFunctionMarketMaker) (*math.Int, error) {
	left, err := CalculateX(y.Add(dy), k, cfmm)
	if err != nil {
		return nil, err
	}

	right, err := CalculateX(y, k, cfmm)
	if err != nil {
		return nil, err
	}

	diff := left.Sub(*right)

	return &diff, nil
}

func CalculateY(x math.Int, k math.LegacyDec, cfmm cfmm.ConstantFunctionMarketMaker) (*math.Int, error) {
	return nil, nil
}

func CalculateDy(x math.Int, dx math.Int, k math.LegacyDec, cfmm cfmm.ConstantFunctionMarketMaker) (*math.Int, error) {
	left, err := CalculateY(x.Add(dx), k, cfmm)
	if err != nil {
		return nil, err
	}

	right, err := CalculateY(x, k, cfmm)
	if err != nil {
		return nil, err
	}

	diff := left.Sub(*right)

	return &diff, nil
}

func CalculateK(x math.Int, y math.Int, cfmm cfmm.ConstantFunctionMarketMaker) (*math.LegacyDec, error) {
	return nil, nil
}

func CalculatePrice(x math.Int, y math.Int, cfmm cfmm.ConstantFunctionMarketMaker) (*math.LegacyDec, error) {
	kValue, err := CalculateK(x, y, cfmm)
	if err != nil {
		return nil, err
	}
	dx := math.NewInt(1000)
	dy, err := CalculateDy(x, dx, *kValue, cfmm)
	if err != nil {
		return nil, err
	}

	price := dy.ToLegacyDec().Quo(dx.ToLegacyDec())

	return &price, nil
}
