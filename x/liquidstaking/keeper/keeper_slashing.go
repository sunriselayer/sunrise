package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stakingtypes "cosmossdk.io/x/staking/types"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k Keeper) GetStakedAmount(ctx context.Context, validatorAddr string) (math.Int, error) {
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	res, err := k.Environment.QueryRouterService.Invoke(ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: moduleAddr.String(),
		ValidatorAddr: validatorAddr,
	})
	if err != nil {
		return math.Int{}, err
	}
	queryDelegationResponse, ok := res.(*stakingtypes.QueryDelegationResponse)
	if !ok {
		return math.Int{}, sdkerrors.ErrInvalidRequest
	}
	stakedAmount := queryDelegationResponse.DelegationResponse.Balance.Amount

	return stakedAmount, nil
}
