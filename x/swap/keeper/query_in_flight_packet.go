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

func (k Keeper) OutgoingInFlightPacketAll(ctx context.Context, req *types.QueryAllOutgoingInFlightPacketRequest) (*types.QueryAllOutgoingInFlightPacketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var inFlightPackets []types.OutgoingInFlightPacket

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	inFlightPacketStore := prefix.NewStore(store, types.KeyPrefix(types.OutgoingInFlightPacketKeyPrefix))

	pageRes, err := query.Paginate(inFlightPacketStore, req.Pagination, func(key []byte, value []byte) error {
		var inFlightPacket types.OutgoingInFlightPacket
		if err := k.cdc.Unmarshal(value, &inFlightPacket); err != nil {
			return err
		}

		inFlightPackets = append(inFlightPackets, inFlightPacket)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOutgoingInFlightPacketResponse{OutgoingInFlightPacket: inFlightPackets, Pagination: pageRes}, nil
}

func (k Keeper) OutgoingInFlightPacket(ctx context.Context, req *types.QueryGetOutgoingInFlightPacketRequest) (*types.QueryGetOutgoingInFlightPacketResponse, error) {
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

	return &types.QueryGetOutgoingInFlightPacketResponse{OutgoingInFlightPacket: val}, nil
}
