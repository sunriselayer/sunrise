package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func (q queryServer) Borrow(ctx context.Context, req *types.QueryBorrowRequest) (*types.QueryBorrowResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	borrow, err := q.k.Borrows.Get(ctx, req.BorrowId)
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "borrow not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBorrowResponse{Borrow: borrow}, nil
}

func (q queryServer) UserBorrows(ctx context.Context, req *types.QueryUserBorrowsRequest) (*types.QueryUserBorrowsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Borrower == "" {
		return nil, status.Error(codes.InvalidArgument, "borrower address cannot be empty")
	}

	borrows, pageResp, err := query.CollectionFilteredPaginate(
		ctx,
		q.k.Borrows,
		req.Pagination,
		func(key uint64, value types.Borrow) (include bool, err error) {
			// Only include borrows for the requested borrower
			return value.Borrower == req.Borrower, nil
		},
		func(key uint64, value types.Borrow) (types.Borrow, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryUserBorrowsResponse{
		Borrows:    borrows,
		Pagination: pageResp,
	}, nil
}