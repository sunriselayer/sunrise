package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "cosmossdk.io/x/staking/types"
)

// StakingKeeper is expected keeper for staking module
type StakingKeeper interface {
	GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error)
	GetDelegatorBonded(ctx context.Context, delegator sdk.AccAddress) (math.Int, error)
}

type TokenConverterKeeper interface {
	DelegateOrSelfDelegate(
		ctx context.Context,
		msg *stakingtypes.MsgDelegate,
		originalFunc func(ctx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error),
	) (*stakingtypes.MsgDelegateResponse, error)
	UndelegateOrSelfUndelegate(
		ctx context.Context,
		msg *stakingtypes.MsgUndelegate,
		originalFunc func(ctx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error),
	) (*stakingtypes.MsgUndelegateResponse, error)
}
