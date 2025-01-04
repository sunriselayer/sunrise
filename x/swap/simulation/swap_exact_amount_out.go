package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"sunrise/x/swap/keeper"
	"sunrise/x/swap/types"
)

func MsgSwapExactAmountOutFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgSwapExactAmountOut] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgSwapExactAmountOut) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgSwapExactAmountOut{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the SwapExactAmountOut simulation

		return []simsx.SimAccount{from}, msg
	}
}
