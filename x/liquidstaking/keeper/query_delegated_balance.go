package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunrise-zone/sunrise-app/x/liquidstaking/types"

	sdkmath "cosmossdk.io/math"
	vestingexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) DelegatedBalance(goCtx context.Context, req *types.QueryDelegatedBalanceRequest) (*types.QueryDelegatedBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	delegator, err := sdk.AccAddressFromBech32(req.Delegator)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid delegator address: %s", err)
	}

	delegated := k.getDelegatedBalance(ctx, delegator)

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}
	vesting := k.getVesting(ctx, delegator).AmountOf(bondDenom)

	vestingDelegated := sdkmath.MinInt(vesting, delegated)
	vestedDelegated := delegated.Sub(vestingDelegated)

	res := types.QueryDelegatedBalanceResponse{
		Vested:  sdk.NewCoin(bondDenom, vestedDelegated),
		Vesting: sdk.NewCoin(bondDenom, vestingDelegated),
	}
	return &res, nil
}

func (k Keeper) getDelegatedBalance(ctx sdk.Context, delegator sdk.AccAddress) sdkmath.Int {
	balance := sdkmath.LegacyZeroDec()

	k.stakingKeeper.IterateDelegatorDelegations(ctx, delegator, func(delegation stakingtypes.Delegation) bool {
		valAddr, _ := sdk.ValAddressFromBech32(delegation.GetValidatorAddr())
		validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
		if err != nil {
			panic(fmt.Sprintf("validator %s for delegation not found", delegation.GetValidatorAddr()))
		}
		tokens := validator.TokensFromSharesTruncated(delegation.GetShares())
		balance = balance.Add(tokens)

		return false
	})
	return balance.TruncateInt()
}

func (k Keeper) getVesting(ctx sdk.Context, delegator sdk.AccAddress) sdk.Coins {
	acc := k.accountKeeper.GetAccount(ctx, delegator)
	if acc == nil {
		// account doesn't exist so amount vesting is 0
		return nil
	}
	vestAcc, ok := acc.(vestingexported.VestingAccount)
	if !ok {
		// account is not vesting type, so amount vesting is 0
		return nil
	}
	return vestAcc.GetVestingCoins(ctx.BlockTime())
}
