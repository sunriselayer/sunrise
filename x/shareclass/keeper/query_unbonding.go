package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func (q queryServer) AddressUnbondings(ctx context.Context, req *types.QueryAddressUnbondingsRequest) (*types.QueryAddressUnbondingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	addr, err := q.k.addressCodec.StringToBytes(req.Address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	unbondings, err := q.k.GetUnbondingsByAddress(ctx, addr)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAddressUnbondingsResponse{Unbondings: unbondings}, nil
}
