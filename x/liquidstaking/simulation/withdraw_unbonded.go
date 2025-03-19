package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/liquidstaking/keeper"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func MsgWithdrawUnbondedFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgWithdrawUnbonded] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgWithdrawUnbonded) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgWithdrawUnbonded{
			Sender: from.AddressBech32,
		}

		// TODO: Handle the WithdrawUnbonded simulation

		return []simsx.SimAccount{from}, msg
	}
}
