package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (q queryServer) Votes(ctx context.Context, req *types.QueryVotesRequest) (*types.QueryVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	votes, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Votes,
		req.Pagination,
		func(key sdk.AccAddress, value types.Vote) (types.Vote, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryVotesResponse{Votes: votes, Pagination: pageRes}, nil
}

func (q queryServer) Vote(ctx context.Context, req *types.QueryVoteRequest) (*types.QueryVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found, err := q.k.GetVote(
		ctx,
		req.Address,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryVoteResponse{Vote: val}, nil
}

func (q queryServer) TallyResult(ctx context.Context, req *types.QueryTallyResultRequest) (*types.QueryTallyResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	gauges, err := q.k.Tally(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTallyResultResponse{Gauges: gauges}, nil
}
