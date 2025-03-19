package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/liquidstaking/keeper"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func MsgLiquidStakeFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgLiquidStake] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgLiquidStake) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgLiquidStake{
			Sender: from.AddressBech32,
		}

		// TODO: Handle the LiquidStake simulation

		return []simsx.SimAccount{from}, msg
	}
}
