package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/swap/types"
)

func (q queryServer) OutgoingInFlightPackets(ctx context.Context, req *types.QueryOutgoingInFlightPacketsRequest) (*types.QueryOutgoingInFlightPacketsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	outgoingInFlightPackets, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.OutgoingInFlightPackets,
		req.Pagination,
		func(key collections.Triple[string, string, uint64], value types.OutgoingInFlightPacket) (types.OutgoingInFlightPacket, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOutgoingInFlightPacketsResponse{Packets: outgoingInFlightPackets, Pagination: pageRes}, nil
}

func (q queryServer) OutgoingInFlightPacket(ctx context.Context, req *types.QueryOutgoingInFlightPacketRequest) (*types.QueryOutgoingInFlightPacketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found, err := q.k.GetOutgoingInFlightPacket(
		ctx,
		req.SrcPortId,
		req.SrcChannelId,
		req.Sequence,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryOutgoingInFlightPacketResponse{Packet: val}, nil
}
