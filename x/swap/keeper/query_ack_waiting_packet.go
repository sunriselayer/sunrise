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

func (k Keeper) IncomingInFlightPacketAll(ctx context.Context, req *types.QueryAllIncomingInFlightPacketRequest) (*types.QueryAllIncomingInFlightPacketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var ackWaitingPackets []types.IncomingInFlightPacket

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	ackWaitingPacketStore := prefix.NewStore(store, types.KeyPrefix(types.IncomingInFlightPacketKeyPrefix))

	pageRes, err := query.Paginate(ackWaitingPacketStore, req.Pagination, func(key []byte, value []byte) error {
		var ackWaitingPacket types.IncomingInFlightPacket
		if err := k.cdc.Unmarshal(value, &ackWaitingPacket); err != nil {
			return err
		}

		ackWaitingPackets = append(ackWaitingPackets, ackWaitingPacket)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllIncomingInFlightPacketResponse{IncomingInFlightPacket: ackWaitingPackets, Pagination: pageRes}, nil
}

func (k Keeper) IncomingInFlightPacket(ctx context.Context, req *types.QueryGetIncomingInFlightPacketRequest) (*types.QueryGetIncomingInFlightPacketResponse, error) {
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

	return &types.QueryGetIncomingInFlightPacketResponse{IncomingInFlightPacket: val}, nil
}
