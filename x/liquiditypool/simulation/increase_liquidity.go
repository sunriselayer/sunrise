package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
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
