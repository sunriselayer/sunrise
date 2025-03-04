package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/swap/types"
)

func (q queryServer) IncomingInFlightPackets(ctx context.Context, req *types.QueryIncomingInFlightPacketsRequest) (*types.QueryIncomingInFlightPacketsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	incomingPackets, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.IncomingInFlightPackets,
		req.Pagination,
		func(key collections.Triple[string, string, uint64], value types.IncomingInFlightPacket) (types.IncomingInFlightPacket, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryIncomingInFlightPacketsResponse{Packets: incomingPackets, Pagination: pageRes}, nil
}

func (q queryServer) IncomingInFlightPacket(ctx context.Context, req *types.QueryIncomingInFlightPacketRequest) (*types.QueryIncomingInFlightPacketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := q.k.GetIncomingInFlightPacket(
		ctx,
		req.SrcPortId, req.SrcChannelId, req.Sequence,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryIncomingInFlightPacketResponse{Packet: val}, nil
}
