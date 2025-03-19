package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/liquidstaking/keeper"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func MsgLiquidUnstakeFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgLiquidUnstake] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgLiquidUnstake) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgLiquidUnstake{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the LiquidUnstake simulation

		return []simsx.SimAccount{from}, msg
	}
}
