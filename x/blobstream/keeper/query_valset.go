package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"sunrise/x/blobstream/types"
)

// TODO add unit tests for all of these requests

// LatestValsetRequestBeforeNonce queries the latest valset request before nonce
func (k Keeper) LatestValsetRequestBeforeNonce(
	c context.Context,
	req *types.QueryLatestValsetRequestBeforeNonceRequest,
) (*types.QueryLatestValsetRequestBeforeNonceResponse, error) {
	vs, err := k.GetLatestValsetBeforeNonce(sdk.UnwrapSDKContext(c), req.Nonce)
	if err != nil {
		return nil, err
	}
	return &types.QueryLatestValsetRequestBeforeNonceResponse{Valset: vs}, nil
}
