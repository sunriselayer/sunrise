package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) Epochs(ctx context.Context, req *types.QueryEpochsRequest) (*types.QueryEpochsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var epochs []types.Epoch

	store := runtime.KVStoreAdapter(q.k.KVStoreService.OpenKVStore(ctx))
	epochStore := prefix.NewStore(store, types.KeyPrefix(types.EpochKey))

	pageRes, err := query.Paginate(epochStore, req.Pagination, func(key []byte, value []byte) error {
		var epoch types.Epoch
		if err := q.k.cdc.Unmarshal(value, &epoch); err != nil {
			return err
		}

		epochs = append(epochs, epoch)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryEpochsResponse{Epochs: epochs, Pagination: pageRes}, nil
}

func (q queryServer) Epoch(ctx context.Context, req *types.QueryEpochRequest) (*types.QueryEpochResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	epoch, found := q.k.GetEpoch(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryEpochResponse{Epoch: epoch}, nil
}
