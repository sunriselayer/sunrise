package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (q queryServer) ClaimableRewards(ctx context.Context, req *types.QueryClaimableRewardsRequest) (*types.QueryClaimableRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	address, err := q.k.addressCodec.StringToBytes(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	validatorAddr, err := q.k.stakingKeeper.ValidatorAddressCodec().StringToBytes(req.ValidatorAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid validator address")
	}

	rewards, err := q.k.GetClaimableRewards(ctx, address, validatorAddr)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryClaimableRewardsResponse{Amount: rewards}, nil
}
