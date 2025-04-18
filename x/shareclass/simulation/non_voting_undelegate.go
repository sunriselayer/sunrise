package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/shareclass/keeper"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func MsgNonVotingUndelegateFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgNonVotingUndelegate] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgNonVotingUndelegate) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgNonVotingUndelegate{
			Sender: from.AddressBech32,
		}

		// TODO: Handle the NonVotingUndelegate simulation

		return []simsx.SimAccount{from}, msg
	}
}
