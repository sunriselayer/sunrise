package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func MsgDecreaseLiquidityFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgDecreaseLiquidity] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgDecreaseLiquidity) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgDecreaseLiquidity{
			Sender: from.AddressBech32,
		}

		// TODO: Handle the DecreaseLiquidity simulation

		return []simsx.SimAccount{from}, msg
	}
}
