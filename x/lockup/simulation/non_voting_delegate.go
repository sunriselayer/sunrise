package simulation

import (
	"context"

	"github.com/cosmos/cosmos-sdk/simsx"

	"github.com/sunriselayer/sunrise/x/lockup/keeper"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func MsgNonVotingDelegateFactory(k keeper.Keeper) simsx.SimMsgFactoryFn[*types.MsgNonVotingDelegate] {
	return func(ctx context.Context, testData *simsx.ChainDataSource, reporter simsx.SimulationReporter) ([]simsx.SimAccount, *types.MsgNonVotingDelegate) {
		from := testData.AnyAccount(reporter)

		msg := &types.MsgNonVotingDelegate{
			Owner: from.AddressBech32,
		}

		// TODO: Handle the NonVotingDelegate simulation

		return []simsx.SimAccount{from}, msg
	}
}
