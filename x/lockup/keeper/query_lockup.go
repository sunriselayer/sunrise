package keeper

import (
	"context"

	"cosmossdk.io/math"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (q queryServer) LockupAccounts(ctx context.Context, req *types.QueryLockupAccountsRequest) (*types.QueryLockupAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	owner, err := q.k.addressCodec.StringToBytes(req.Owner)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	lockupAccounts, err := q.k.GetLockupAccountsByOwner(ctx, owner)
	if err != nil {
		return nil, err
	}

	return &types.QueryLockupAccountsResponse{
		LockupAccounts: lockupAccounts,
	}, nil
}

func (q queryServer) LockupAccount(ctx context.Context, req *types.QueryLockupAccountRequest) (*types.QueryLockupAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	owner, err := q.k.addressCodec.StringToBytes(req.Owner)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	lockupAccount, err := q.k.GetLockupAccount(ctx, owner, req.LockupAccountId)
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

func (q queryServer) SpendableAmount(ctx context.Context, req *types.QuerySpendableAmountRequest) (*types.QuerySpendableAmountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	owner, err := q.k.addressCodec.StringToBytes(req.Owner)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	lockupAccount, err := q.k.GetLockupAccount(ctx, owner, req.LockupAccountId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	lockupAcc, err := q.k.addressCodec.StringToBytes(lockupAccount.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	transferableDenom, err := q.k.tokenConverterKeeper.GetTransferableDenom(ctx)
	if err != nil {
		return nil, err
	}

	balance := q.k.bankKeeper.GetBalance(ctx, lockupAcc, transferableDenom)

	currentTime := sdk.UnwrapSDKContext(ctx).BlockTime().Unix()
	_, lockedAmount, err := lockupAccount.GetLockCoinInfo(currentTime)
	if err != nil {
		return nil, err
	}

	// refresh ubd entries to make sure delegation locking amount is up to date
	err = q.k.CheckUnbondingEntriesMature(ctx, owner, req.LockupAccountId)
	if err != nil {
		return nil, err
	}

	notBondedLockedAmount := lockupAccount.GetNotBondedLockedAmount(lockedAmount)
	spendable, err := balance.Amount.SafeSub(notBondedLockedAmount)
	if err != nil {
		spendable = math.ZeroInt()
	}

	return &types.QuerySpendableAmountResponse{
		SpendableAmount: spendable.String(),
	}, nil
}
