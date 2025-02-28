package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) PublishedData(goCtx context.Context, req *types.QueryPublishedDataRequest) (*types.QueryPublishedDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	data, found, err := q.k.GetPublishedData(ctx, req.MetadataUri)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, types.ErrDataNotFound
	}
	return &types.QueryPublishedDataResponse{Data: data}, nil
}

func (q queryServer) AllPublishedData(goCtx context.Context, req *types.QueryAllPublishedDataRequest) (*types.QueryAllPublishedDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAllPublishedDataResponse{Data: q.k.GetAllPublishedData(ctx)}, nil
}

func (q queryServer) ZkpProofThreshold(goCtx context.Context, req *types.QueryZkpProofThresholdRequest) (*types.QueryZkpProofThresholdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryZkpProofThresholdResponse{Threshold: q.k.GetZkpThreshold(ctx, req.ShardCount)}, nil
}

func (q queryServer) ProofDeputy(goCtx context.Context, req *types.QueryProofDeputyRequest) (*types.QueryProofDeputyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	validator, err := q.k.validatorAddressCodec.StringToBytes(req.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}
	deputy, found, err := q.k.GetProofDeputy(ctx, validator)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, types.ErrDeputyNotFound
	}

	return &types.QueryProofDeputyResponse{DeputyAddress: sdk.AccAddress(deputy).String()}, nil
}

func (q queryServer) ValidityProof(goCtx context.Context, req *types.QueryValidityProofRequest) (*types.QueryValidityProofResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	validator, err := q.k.validatorAddressCodec.StringToBytes(req.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}
	proof, found, err := q.k.GetProof(ctx, req.MetadataUri, validator)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, types.ErrProofNotFound
	}

	return &types.QueryValidityProofResponse{Proof: proof}, nil
}

func (q queryServer) AllValidityProofs(goCtx context.Context, req *types.QueryAllValidityProofsRequest) (*types.QueryAllValidityProofsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	proofs, err := q.k.GetAllProofs(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllValidityProofsResponse{Proofs: proofs}, nil
}

func (q queryServer) Invalidity(goCtx context.Context, req *types.QueryInvalidityRequest) (*types.QueryInvalidityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := q.k.addressCodec.StringToBytes(req.SenderAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}
	invalidity, found, err := q.k.GetInvalidity(ctx, req.MetadataUri, sender)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, types.ErrProofNotFound
	}

	return &types.QueryInvalidityResponse{Invalidity: invalidity}, nil
}

func (q queryServer) AllInvalidity(goCtx context.Context, req *types.QueryAllInvalidityRequest) (*types.QueryAllInvalidityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	invalidities, err := q.k.GetAllInvalidities(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllInvalidityResponse{Invalidity: invalidities}, nil
}
