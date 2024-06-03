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

func (k Keeper) InFlightPacketAll(ctx context.Context, req *types.QueryAllInFlightPacketRequest) (*types.QueryAllInFlightPacketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var inFlightPackets []types.InFlightPacket

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	inFlightPacketStore := prefix.NewStore(store, types.KeyPrefix(types.InFlightPacketKeyPrefix))

	pageRes, err := query.Paginate(inFlightPacketStore, req.Pagination, func(key []byte, value []byte) error {
		var inFlightPacket types.InFlightPacket
		if err := k.cdc.Unmarshal(value, &inFlightPacket); err != nil {
			return err
		}

		inFlightPackets = append(inFlightPackets, inFlightPacket)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllInFlightPacketResponse{InFlightPacket: inFlightPackets, Pagination: pageRes}, nil
}

func (k Keeper) InFlightPacket(ctx context.Context, req *types.QueryGetInFlightPacketRequest) (*types.QueryGetInFlightPacketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetInFlightPacket(
		ctx,
		req.SrcPortId,
		req.SrcChannelId,
		req.Sequence,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetInFlightPacketResponse{InFlightPacket: val}, nil
}
