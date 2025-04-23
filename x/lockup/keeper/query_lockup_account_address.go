package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (q queryServer) LockupAccount(ctx context.Context, req *types.QueryLockupAccountRequest) (*types.QueryLockupAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	owner, err := q.k.addressCodec.StringToBytes(req.Owner)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	lockupAccount, err := q.k.GetLockupAccount(ctx, owner, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	currentTime := sdk.UnwrapSDKContext(ctx).BlockTime().Unix()
	unlockedAmount, lockedAmount, err := lockupAccount.GetLockCoinInfo(currentTime)
	if err != nil {
		return nil, err
	}

	return &types.QueryLockupAccountResponse{
		LockupAccount:  lockupAccount,
		UnlockedAmount: unlockedAmount.String(),
		LockedAmount:   lockedAmount.String(),
	}, nil
}
