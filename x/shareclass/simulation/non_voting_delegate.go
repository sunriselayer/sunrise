package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/shareclass/keeper"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

func MsgNonVotingDelegateFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgNonVotingDelegate] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgNonVotingDelegate) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgNonVotingDelegate{
			Creator: from.AddressBech32,
		}

		// TODO: Handle the NonVotingDelegate simulation

		return []simsx.SimAccount{from}, msg
	}
}
