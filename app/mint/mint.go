package mint

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/math"
	authtypes "cosmossdk.io/x/auth/types"
	banktypes "cosmossdk.io/x/bank/types"
	minttypes "cosmossdk.io/x/mint/types"
	stakingtypes "cosmossdk.io/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface{}

func ProvideMintFn(bankKeeper BankKeeper) minttypes.MintFn {
	return func(ctx context.Context, env appmodule.Environment, minter *minttypes.Minter, epochID string, epochNumber int64) error {
		return nil
	}
}
