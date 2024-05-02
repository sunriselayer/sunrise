package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sunriselayer/sunrise-app/x/blobgrant/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RegistrationAll(ctx context.Context, req *types.QueryAllRegistrationRequest) (*types.QueryAllRegistrationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var registrations []types.Registration

    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	registrationStore := prefix.NewStore(store, types.KeyPrefix(types.RegistrationKeyPrefix))

	pageRes, err := query.Paginate(registrationStore, req.Pagination, func(key []byte, value []byte) error {
		var registration types.Registration
		if err := k.cdc.Unmarshal(value, &registration); err != nil {
			return err
		}

		registrations = append(registrations, registration)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRegistrationResponse{Registration: registrations, Pagination: pageRes}, nil
}

func (k Keeper) Registration(ctx context.Context, req *types.QueryGetRegistrationRequest) (*types.QueryGetRegistrationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetRegistration(
	    ctx,
	    req.Address,
        )
	if !found {
	    return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetRegistrationResponse{Registration: val}, nil
}