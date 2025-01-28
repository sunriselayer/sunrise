package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/selfdelegation/types"
)

func (q queryServer) QuerySelfDelegationProxyAccountByOwner(ctx context.Context, req *types.QuerySelfDelegationProxyAccountByOwnerRequest) (*types.QuerySelfDelegationProxyAccountByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ownerAddress, err := q.k.addressCodec.StringToBytes(req.OwnerAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid owner address")
	}

	proxyAccountAddress, err := q.k.SelfDelegationProxies.Get(ctx, ownerAddress)
	if err != nil {
		return nil, status.Error(codes.NotFound, "proxy account not found")
	}

	proxyAccountString, err := q.k.addressCodec.BytesToString(proxyAccountAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert proxy account address to string")
	}

	return &types.QuerySelfDelegationProxyAccountByOwnerResponse{SelfDelegationProxyAccountAddress: proxyAccountString}, nil
}
