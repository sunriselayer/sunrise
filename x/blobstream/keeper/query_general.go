package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunrise-zone/sunrise-app/x/blobstream/types"
)

// LatestUnbondingHeight queries the latest unbonding height.
func (k Keeper) LatestUnbondingHeight(
	c context.Context,
	_ *types.QueryLatestUnbondingHeightRequest,
) (*types.QueryLatestUnbondingHeightResponse, error) {
	return &types.QueryLatestUnbondingHeightResponse{
		Height: k.GetLatestUnBondingBlockHeight(sdk.UnwrapSDKContext(c)),
	}, nil
}

// EarliestAttestationNonce queries the earliest attestation nonce.
func (k Keeper) EarliestAttestationNonce(
	c context.Context,
	_ *types.QueryEarliestAttestationNonceRequest,
) (*types.QueryEarliestAttestationNonceResponse, error) {
	return &types.QueryEarliestAttestationNonceResponse{
		Nonce: k.GetEarliestAvailableAttestationNonce(sdk.UnwrapSDKContext(c)),
	}, nil
}

// EVMAddress tries to find the associated EVM address for a given validator address. If
// none is found, an empty address is returned
func (k Keeper) EvmAddress(goCtx context.Context, req *types.QueryEvmAddressRequest) (*types.QueryEvmAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	evmAddr, exists := k.GetEVMAddress(ctx, valAddr)
	if !exists {
		return &types.QueryEvmAddressResponse{}, nil
	}
	return &types.QueryEvmAddressResponse{
		EvmAddress: evmAddr.Hex(),
	}, nil
}
