package keeper

import (
	"bytes"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/selfdelegation/types"
)

func (q queryServer) LockupAccountsByOwner(ctx context.Context, req *types.QueryLockupAccountsByOwnerRequest) (*types.QueryLockupAccountsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ownerAddress, err := q.k.addressCodec.StringToBytes(req.OwnerAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid owner address")
	}

	var lockupAccountAddresses []string
	err = q.k.LockupAccounts.Walk(ctx, nil, func(key []byte, val []byte) (stop bool, err error) {
		if bytes.Equal(val, ownerAddress) { // if owner address matches
			lockupAddressString, err := q.k.addressCodec.BytesToString(key)
			if err != nil {
				// continue walking and skip current lockup account
				return true, nil
			}
			lockupAccountAddresses = append(lockupAccountAddresses, lockupAddressString)
		}
		return false, nil // continue walking to find all lockup accounts for the owner
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to walk lockup accounts")
	}

	return &types.QueryLockupAccountsByOwnerResponse{LockupAccountAddresses: lockupAccountAddresses}, nil
}
