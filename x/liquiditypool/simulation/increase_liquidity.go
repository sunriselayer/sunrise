package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"sunrise/x/liquiditypool/keeper"
	"sunrise/x/liquiditypool/types"
)

func MsgIncreaseLiquidityFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgIncreaseLiquidity] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgIncreaseLiquidity) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgIncreaseLiquidity{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the IncreaseLiquidity simulation

		return []simsx.SimAccount{from}, msg
	}
}
