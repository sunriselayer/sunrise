package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k Keeper) WrapPositionInfo(ctx context.Context, position types.Position) types.PositionInfo {
	pool, found, err := k.GetPool(ctx, position.PoolId)
	if err != nil {
		return types.PositionInfo{Position: position}
	}
	if !found {
		return types.PositionInfo{Position: position}
	}

	liquidity, err := math.LegacyNewDecFromStr(position.Liquidity)
	if err != nil {
		return types.PositionInfo{Position: position}
	}
	actualAmountBase, actualAmountQuote, err := pool.CalcActualAmounts(position.LowerTick, position.UpperTick, liquidity)
	if err != nil {
		return types.PositionInfo{Position: position}
	}

	return types.PositionInfo{
		Position:   position,
		TokenBase:  sdk.NewCoin(pool.DenomBase, actualAmountBase.TruncateInt()),
		TokenQuote: sdk.NewCoin(pool.DenomQuote, actualAmountQuote.TruncateInt()),
	}
}

func (q queryServer) Positions(ctx context.Context, req *types.QueryPositionsRequest) (*types.QueryPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	positions, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Positions,
		req.Pagination,
		func(_ uint64, value types.Position) (types.PositionInfo, error) {
			return q.k.WrapPositionInfo(ctx, value), nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPositionsResponse{Positions: positions, Pagination: pageRes}, nil
}

func (q queryServer) Position(ctx context.Context, req *types.QueryPositionRequest) (*types.QueryPositionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	position, found := q.k.GetPosition(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryPositionResponse{Position: q.k.WrapPositionInfo(ctx, position)}, nil
}

func (q queryServer) PoolPositions(ctx context.Context, req *types.QueryPoolPositionsRequest) (*types.QueryPoolPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	positionInfos := []types.PositionInfo{}
	positions := q.k.GetPositionsByPool(ctx, req.PoolId)
	for _, position := range positions {
		positionInfos = append(positionInfos, q.k.WrapPositionInfo(ctx, position))
	}

	return &types.QueryPoolPositionsResponse{Positions: positionInfos}, nil
}

func (q queryServer) AddressPositions(ctx context.Context, req *types.QueryAddressPositionsRequest) (*types.QueryAddressPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	addr, err := q.k.addressCodec.StringToBytes(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	positionInfos := []types.PositionInfo{}
	positions := q.k.GetPositionsByAddress(ctx, addr)
	for _, position := range positions {
		positionInfos = append(positionInfos, q.k.WrapPositionInfo(ctx, position))
	}

	return &types.QueryAddressPositionsResponse{Positions: positionInfos}, nil
}

func (q queryServer) PositionFees(ctx context.Context, req *types.QueryPositionFeesRequest) (*types.QueryPositionFeesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectedFees, err := q.k.GetClaimableFees(sdk.UnwrapSDKContext(ctx), req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryPositionFeesResponse{Fees: collectedFees}, nil
}
