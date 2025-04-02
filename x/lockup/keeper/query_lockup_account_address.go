package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (q queryServer) LockupAccountAddress(ctx context.Context, req *types.QueryLockupAccountAddressRequest) (*types.QueryLockupAccountAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryLockupAccountAddressResponse{}, nil
}
