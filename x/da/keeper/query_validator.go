package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func (q queryServer) ValidatorsPowerSnapshot(goCtx context.Context, req *types.QueryValidatorsPowerSnapshotRequest) (*types.QueryValidatorsPowerSnapshotResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	validatorsPowerSnapshot, found, err := q.k.GetValidatorsPowerSnapshot(ctx, req.BlockHeight)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Error(codes.NotFound, "validators power snapshot not found")
	}

	return &types.QueryValidatorsPowerSnapshotResponse{ValidatorsPowerSnapshot: validatorsPowerSnapshot}, nil
}

func (q queryServer) CommitmentKey(goCtx context.Context, req *types.QueryCommitmentKeyRequest) (*types.QueryCommitmentKeyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	validatorAddress, err := q.k.addressCodec.StringToBytes(req.ValidatorAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid validator address")
	}

	data, found, err := q.k.GetCommitmentKey(ctx, validatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Error(codes.NotFound, "commitment key not found")
	}

	return &types.QueryCommitmentKeyResponse{Pubkey: data.Pubkey}, nil
}

func (q queryServer) ShardIndices(goCtx context.Context, req *types.QueryShardIndicesRequest) (*types.QueryShardIndicesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	validatorAddress, err := q.k.addressCodec.StringToBytes(req.ValidatorAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid validator address")
	}

	indices, err := types.CorrespondingShardIndices(req.ShardsMerkleRoot, validatorAddress, req.ShardCount, req.ShardCountPerValidator)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	indicesArray := []uint32{}
	for index := range indices {
		indicesArray = append(indicesArray, index)
	}

	return &types.QueryShardIndicesResponse{ShardIndices: indicesArray}, nil
}
