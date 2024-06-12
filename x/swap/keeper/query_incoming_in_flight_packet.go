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

func (k Keeper) IncomingInFlightPackets(ctx context.Context, req *types.QueryIncomingInFlightPacketsRequest) (*types.QueryIncomingInFlightPacketsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var incomingPackets []types.IncomingInFlightPacket

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	incomingPacketStore := prefix.NewStore(store, types.KeyPrefix(types.IncomingInFlightPacketKeyPrefix))

	pageRes, err := query.Paginate(incomingPacketStore, req.Pagination, func(key []byte, value []byte) error {
		var incomingPacket types.IncomingInFlightPacket
		if err := k.cdc.Unmarshal(value, &incomingPacket); err != nil {
			return err
		}

		incomingPackets = append(incomingPackets, incomingPacket)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryIncomingInFlightPacketsResponse{Packets: incomingPackets, Pagination: pageRes}, nil
}

func (k Keeper) IncomingInFlightPacket(ctx context.Context, req *types.QueryIncomingInFlightPacketRequest) (*types.QueryIncomingInFlightPacketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetIncomingInFlightPacket(
		ctx,
		req.SrcPortId, req.SrcChannelId, req.Sequence,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryIncomingInFlightPacketResponse{Packet: val}, nil
}
