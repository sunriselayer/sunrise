package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) WrapPositionInfo(ctx context.Context, position types.Position) types.PositionInfo {
	pool, found := k.GetPool(ctx, position.PoolId)
	if !found {
		return types.PositionInfo{Position: position}
	}

	actualAmountBase, actualAmountQuote, err := pool.CalcActualAmounts(position.LowerTick, position.UpperTick, position.Liquidity)
	if err != nil {
		return types.PositionInfo{Position: position}
	}

	return types.PositionInfo{
		Position:   position,
		TokenBase:  sdk.NewCoin(pool.DenomBase, actualAmountBase.TruncateInt()),
		TokenQuote: sdk.NewCoin(pool.DenomQuote, actualAmountQuote.TruncateInt()),
	}
}

func (k Keeper) Positions(ctx context.Context, req *types.QueryPositionsRequest) (*types.QueryPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var positions []types.PositionInfo

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	positionStore := prefix.NewStore(store, types.KeyPrefix(types.PositionKey))

	pageRes, err := query.Paginate(positionStore, req.Pagination, func(key []byte, value []byte) error {
		var position types.Position
		if err := k.cdc.Unmarshal(value, &position); err != nil {
			return err
		}

		positions = append(positions, k.WrapPositionInfo(ctx, position))
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPositionsResponse{Positions: positions, Pagination: pageRes}, nil
}

func (k Keeper) Position(ctx context.Context, req *types.QueryPositionRequest) (*types.QueryPositionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	position, found := k.GetPosition(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryPositionResponse{Position: k.WrapPositionInfo(ctx, position)}, nil
}

func (k Keeper) PoolPositions(ctx context.Context, req *types.QueryPoolPositionsRequest) (*types.QueryPoolPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	positionInfos := []types.PositionInfo{}
	positions := k.GetPositionsByPool(ctx, req.PoolId)
	for _, position := range positions {
		positionInfos = append(positionInfos, k.WrapPositionInfo(ctx, position))
	}

	return &types.QueryPoolPositionsResponse{Positions: positionInfos}, nil
}

func (k Keeper) AddressPositions(ctx context.Context, req *types.QueryAddressPositionsRequest) (*types.QueryAddressPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	positionInfos := []types.PositionInfo{}
	positions := k.GetPositionsByAddress(ctx, req.Address)
	for _, position := range positions {
		positionInfos = append(positionInfos, k.WrapPositionInfo(ctx, position))
	}

	return &types.QueryAddressPositionsResponse{Positions: positionInfos}, nil
}

func (k Keeper) PositionFees(ctx context.Context, req *types.QueryPositionFeesRequest) (*types.QueryPositionFeesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// position, found := k.GetPosition(ctx, req.Id)
	// if !found {
	// 	return nil, sdkerrors.ErrKeyNotFound
	// }
	// k.GetFeeAccumulator(c)
	// TODO:

	return &types.QueryPositionFeesResponse{Fees: sdk.Coins{}}, nil
}
