package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/lockup/keeper"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func MsgNonVotingUndelegateFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgNonVotingUndelegate] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgNonVotingUndelegate) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgNonVotingUndelegate{
			Owner: from.AddressBech32,
		}

		// TODO: Handle the NonVotingUndelegate simulation

		return []simsx.SimAccount{from}, msg
	}
}
