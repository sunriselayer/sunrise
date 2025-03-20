package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "cosmossdk.io/x/staking/types"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k Keeper) ConvertAndDelegate(ctx context.Context, sender sdk.AccAddress, validatorAddr string, amount math.Int) error {
	// Prepare fee and bond coin
	params, err := k.tokenConverterKeeper.GetParams(ctx)
	if err != nil {
		return err
	}
	feeCoin := sdk.NewCoin(params.FeeDenom, amount)
	bondCoin := sdk.NewCoin(params.BondDenom, amount)

	// Send fee coin to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return err
	}

	// Convert fee denom to bond denom
	err = k.tokenConverterKeeper.ConvertReverse(ctx, amount, sender)
	if err != nil {
		return err
	}

	// Stake
	_, err = k.Environment.MsgRouterService.Invoke(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sender.String(),
		ValidatorAddress: validatorAddr,
		Amount:           bondCoin,
	})
	if err != nil {
		return err
	}

	return nil
}
