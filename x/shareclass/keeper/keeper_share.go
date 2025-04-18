package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (k Keeper) GetShare(ctx context.Context, address sdk.AccAddress, validatorAddr string) math.Int {
	shareDenom := types.NonVotingShareTokenDenom(validatorAddr)

	shareSupply := k.bankKeeper.GetBalance(ctx, address, shareDenom)

	return shareSupply.Amount
}

func (k Keeper) GetTotalShare(ctx context.Context, validatorAddr string) math.Int {
	shareDenom := types.NonVotingShareTokenDenom(validatorAddr)

	shareSupply := k.bankKeeper.GetSupply(ctx, shareDenom)

	return shareSupply.Amount
}

func (k Keeper) CalculateAmountByShare(ctx context.Context, validatorAddr string, share math.Int) (math.Int, error) {
	totalShare := k.GetTotalShare(ctx, validatorAddr)
	totalStaked, err := k.GetTotalStakedAmount(ctx, validatorAddr)
	if err != nil {
		return math.Int{}, err
	}
	amount, err := types.CalculateAmountByShare(totalShare, totalStaked, share)
	if err != nil {
		return math.Int{}, err
	}

	return amount, nil
}

func (k Keeper) CalculateShareByAmount(ctx context.Context, validatorAddr string, amount math.Int) (math.Int, error) {
	totalShare := k.GetTotalShare(ctx, validatorAddr)
	// If total share is zero, return the amount
	// GetTotalStakedAmount will result in rpc error: code = NotFound
	if totalShare.IsZero() {
		return amount, nil
	}
	totalStaked, err := k.GetTotalStakedAmount(ctx, validatorAddr)
	if err != nil {
		return math.Int{}, err
	}
	share, err := types.CalculateShareByAmount(totalShare, totalStaked, amount)
	if err != nil {
		return math.Int{}, err
	}

	return share, nil
}
