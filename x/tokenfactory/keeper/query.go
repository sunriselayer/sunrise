package keeper

import (
	"context"

	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenfactory/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

func (q queryServer) DenomAuthorityMetadata(ctx context.Context, req *types.QueryDenomAuthorityMetadataRequest) (*types.QueryDenomAuthorityMetadataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	decodedDenom, err := url.QueryUnescape(req.Denom)
	if err == nil {
		req.Denom = decodedDenom
	}
	authorityMetadata, err := q.k.GetAuthorityMetadata(sdkCtx, req.GetDenom())
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomAuthorityMetadataResponse{AuthorityMetadata: authorityMetadata}, nil
}

func (q queryServer) DenomsFromCreator(ctx context.Context, req *types.QueryDenomsFromCreatorRequest) (*types.QueryDenomsFromCreatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	denoms := q.k.getDenomsFromCreator(sdkCtx, req.GetCreator())
	return &types.QueryDenomsFromCreatorResponse{Denoms: denoms}, nil
}

func (q queryServer) BeforeSendHookAddress(ctx context.Context, req *types.QueryBeforeSendHookAddressRequest) (*types.QueryBeforeSendHookAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	decodedDenom, err := url.QueryUnescape(req.Denom)
	if err == nil {
		req.Denom = decodedDenom
	}

	cosmwasmAddress := q.k.GetBeforeSendHook(sdkCtx, req.GetDenom())

	return &types.QueryBeforeSendHookAddressResponse{CosmwasmAddress: cosmwasmAddress}, nil
}

func (q queryServer) AllBeforeSendHooksAddresses(ctx context.Context, req *types.QueryAllBeforeSendHooksAddressesRequest) (*types.QueryAllBeforeSendHooksAddressesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	denoms, beforesendHookAddresses := q.k.GetAllBeforeSendHooks(sdkCtx)

	return &types.QueryAllBeforeSendHooksAddressesResponse{Denoms: denoms, BeforeSendHookAddresses: beforesendHookAddresses}, nil
}
