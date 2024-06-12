package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sunriselayer/sunrise/x/swap/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OutgoingInFlightPackets(ctx context.Context, req *types.QueryOutgoingInFlightPacketsRequest) (*types.QueryOutgoingInFlightPacketsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var outgoingInFlightPackets []types.OutgoingInFlightPacket

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	outgoingInFlightPacketStore := prefix.NewStore(store, types.KeyPrefix(types.OutgoingInFlightPacketKeyPrefix))

	pageRes, err := query.Paginate(outgoingInFlightPacketStore, req.Pagination, func(key []byte, value []byte) error {
		var outgoingInFlightPacket types.OutgoingInFlightPacket
		if err := k.cdc.Unmarshal(value, &outgoingInFlightPacket); err != nil {
			return err
		}

		outgoingInFlightPackets = append(outgoingInFlightPackets, outgoingInFlightPacket)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOutgoingInFlightPacketsResponse{Packets: outgoingInFlightPackets, Pagination: pageRes}, nil
}

func (k Keeper) OutgoingInFlightPacket(ctx context.Context, req *types.QueryOutgoingInFlightPacketRequest) (*types.QueryOutgoingInFlightPacketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetOutgoingInFlightPacket(
		ctx,
		req.SrcPortId,
		req.SrcChannelId,
		req.Sequence,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryOutgoingInFlightPacketResponse{Packet: val}, nil
}
