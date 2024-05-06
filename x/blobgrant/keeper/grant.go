package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/blobgrant/types"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k Keeper) EndBlock(ctx context.Context) error {
	params := k.GetParams(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if sdkCtx.BlockHeight()%int64(params.BlockHeightDuration) != 0 {
		return nil
	}

	// charge balances
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalances := k.bankKeeper.SpendableCoins(ctx, moduleAddr)
	grantAmount := moduleBalances.AmountOf(types.GrantTokenDenom)

	if grantAmount.LT(params.GrantTokenRefillThreshold) {
		mintAmount := params.GrantTokenRefillThreshold.Sub(grantAmount)
		mintCoins := sdk.NewCoins(sdk.NewCoin(types.GrantTokenDenom, mintAmount))
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
			return err
		}
	}

	// iterate registrations
	registrations := k.GetAllRegistration(ctx)

	for _, registration := range registrations {
		addr := sdk.MustAccAddressFromBech32(registration.LiquidityProvider)
		balances := k.bankKeeper.SpendableCoins(ctx, addr)

		for _, balance := range balances {
			var moduleName string
			var poolId uint64
			_, err := fmt.Sscanf(balance.Denom, "%s/%d", &moduleName, &poolId)

			if err != nil || moduleName != liquiditypooltypes.ModuleName {
				continue
			}

			// TODO: calculate SR based value of the LP token
			grantAmount := math.NewInt(0)
			expiry := sdk.UnwrapSDKContext(ctx).BlockTime().Add(params.ExpiryDuration)
			allowance := types.NewFeeAllowance(grantAmount, expiry)

			err = k.feeGrantKeeper.GrantAllowance(ctx, moduleAddr, sdk.AccAddress(registration.Grantee), allowance)
			if err != nil {
				continue
			}
		}
	}

	return nil
}
