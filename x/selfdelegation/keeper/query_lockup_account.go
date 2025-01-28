package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/selfdelegation/types"
)

func (q queryServer) LockupAccountByOwner(ctx context.Context, req *types.QueryLockupAccountByOwnerRequest) (*types.QueryLockupAccountByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ownerAddress, err := q.k.addressCodec.StringToBytes(req.OwnerAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid owner address")
	}

	lockupAddress, err := q.k.LockupAccounts.Get(ctx, ownerAddress)
	if err != nil {
		return nil, status.Error(codes.NotFound, "lockup account not found")
	}

	LockupAddressString, err := q.k.addressCodec.BytesToString(lockupAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert lockup account address to string")
	}

	return &types.QueryLockupAccountByOwnerResponse{LockupAccountAddress: LockupAddressString}, nil
}
