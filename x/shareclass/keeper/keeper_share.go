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
	outputAmount, err := types.CalculateUndelegationOutputAmount(share, totalShare, totalStaked)
	if err != nil {
		return math.Int{}, err
	}

	return outputAmount, nil
}
