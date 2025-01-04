package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"sunrise/x/swap/keeper"
	"sunrise/x/swap/types"
)

func MsgSwapExactAmountInFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgSwapExactAmountIn] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgSwapExactAmountIn) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgSwapExactAmountIn{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the SwapExactAmountIn simulation

		return []simsx.SimAccount{from}, msg
	}
}
