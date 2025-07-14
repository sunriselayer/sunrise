package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (q queryServer) AddressBonded(ctx context.Context, req *types.QueryAddressBondedRequest) (*types.QueryAddressBondedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	address, err := q.k.addressCodec.StringToBytes(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	balances := q.k.bankKeeper.GetAllBalances(ctx, address)

	transferableDenom, err := q.k.tokenConverterKeeper.GetTransferableDenom(ctx)
	if err != nil {
		return nil, err
	}

	bonds := []types.ValidatorBond{}

	for _, balance := range balances {
		matches := types.NonVotingShareTokenDenomRegexp().FindStringSubmatch(balance.Denom)
		if len(matches) > 1 {
			validatorAddr := matches[1]
			bondAmount, err := q.k.CalculateAmountByShare(ctx, validatorAddr, balance.Amount)
			if err != nil {
				return nil, err
			}
			bonds = append(bonds, types.ValidatorBond{
				ValidatorAddress: validatorAddr,
				Amount:           sdk.NewCoin(transferableDenom, bondAmount),
			})
		}
	}

	return &types.QueryAddressBondedResponse{Bonds: bonds}, nil
}
