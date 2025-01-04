package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"sunrise/x/liquiditypool/keeper"
	"sunrise/x/liquiditypool/types"
)

func MsgDecreaseLiquidityFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgDecreaseLiquidity] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgDecreaseLiquidity) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgDecreaseLiquidity{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the DecreaseLiquidity simulation

		return []simsx.SimAccount{from}, msg
	}
}
