package customtypes

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BeforeSendHook func(ctx context.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) error
