package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (q queryServer) Unbondings(ctx context.Context, req *types.QueryUnbondingsRequest) (*types.QueryUnbondingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// unbondings, pageRes, err := query.CollectionPaginate(
	// 	ctx,
	// 	q.k.Unbondings.Indexes.Address,
	// 	req.Pagination,
	// 	func(address sdk.AccAddress, id uint64) (types.Unbonding, error) {
	// 		return value, nil
	// 	},
	// )

	unbondings := []types.Unbonding{}
	err := q.k.Unbondings.Indexes.Address.Walk(ctx, nil,
		func(address sdk.AccAddress, id uint64) (stop bool, err error) {
			unbonding, _, err := q.k.GetUnbonding(ctx, id)
			if err != nil {
				return true, err
			}

			unbondings = append(unbondings, unbonding)
			return false, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryUnbondingsResponse{Unbondings: unbondings}, nil
}
