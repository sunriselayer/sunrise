package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (q queryServer) BlobDeclaration(ctx context.Context, req *types.QueryBlobDeclarationRequest) (*types.QueryBlobDeclarationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	data, found, err := q.k.GetBlobDeclaration(ctx, req.ShardsMerkleRoot)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Error(codes.NotFound, "blob declaration not found")
	}
	return &types.QueryBlobDeclarationResponse{BlobDeclaration: data}, nil
}

func (q queryServer) BlobCommitment(ctx context.Context, req *types.QueryBlobCommitmentRequest) (*types.QueryBlobCommitmentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	data, found, err := q.k.GetBlobCommitment(ctx, req.ShardsMerkleRoot)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Error(codes.NotFound, "blob commitment not found")
	}
	return &types.QueryBlobCommitmentResponse{BlobCommitment: data}, nil
}

func (q queryServer) Challenges(ctx context.Context, req *types.QueryChallengesRequest) (*types.QueryChallengesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	challenges, err := q.k.GetAllChallenges(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryChallengesResponse{Challenges: challenges}, nil
}

func (q queryServer) Challenge(ctx context.Context, req *types.QueryChallengeRequest) (*types.QueryChallengeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	challenge, found, err := q.k.GetChallenge(ctx, req.ShardsMerkleRoot, req.ShardIndex, req.EvaluationPointIndex)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Error(codes.NotFound, "challenge not found")
	}

	return &types.QueryChallengeResponse{Challenge: challenge}, nil
}
